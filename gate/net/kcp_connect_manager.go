package net

import (
	"bytes"
	"context"
	"encoding/binary"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"hk4e/common/config"
	"hk4e/common/mq"
	"hk4e/common/region"
	"hk4e/common/rpc"
	"hk4e/gate/client_proto"
	"hk4e/gate/kcp"
	"hk4e/node/api"
	"hk4e/pkg/logger"
	"hk4e/pkg/random"
	"hk4e/protocol/cmd"
)

const (
	ConnSynPacketFreqLimit = 100        // 连接建立握手包每秒发包频率限制
	RecvPacketFreqLimit    = 1000       // 客户端上行每秒发包频率限制
	SendPacketFreqLimit    = 1000       // 服务器下行每秒发包频率限制
	PacketMaxLen           = 343 * 1024 // 最大应用层包长度
	ConnRecvTimeout        = 30         // 收包超时时间 秒
	ConnSendTimeout        = 10         // 发包超时时间 秒
	MaxClientConnNumLimit  = 1000       // 最大客户端连接数限制
)

var CLIENT_CONN_NUM int32 = 0 // 当前客户端连接数

type KcpConnectManager struct {
	discovery *rpc.DiscoveryClient // node服务器客户端
	openState bool                 // 网关开放状态
	// 会话
	sessionConvIdMap      map[uint64]*Session
	sessionUserIdMap      map[uint32]*Session
	sessionMapLock        sync.RWMutex
	createSessionChan     chan *Session
	destroySessionChan    chan *Session
	globalGsOnlineMap     map[uint32]string
	globalGsOnlineMapLock sync.RWMutex
	// 连接事件
	kcpEventInput            chan *KcpEvent
	kcpEventOutput           chan *KcpEvent
	reLoginRemoteKickRegChan chan *RemoteKick
	// 协议
	serverCmdProtoMap *cmd.CmdProtoMap
	clientCmdProtoMap *client_proto.ClientCmdProtoMap
	// 输入输出管道
	messageQueue *mq.MessageQueue
	// 密钥
	dispatchKey  []byte
	signRsaKey   []byte
	encRsaKeyMap map[string][]byte
}

func NewKcpConnectManager(messageQueue *mq.MessageQueue, discovery *rpc.DiscoveryClient) (r *KcpConnectManager) {
	r = new(KcpConnectManager)
	r.discovery = discovery
	r.openState = true
	r.sessionConvIdMap = make(map[uint64]*Session)
	r.sessionUserIdMap = make(map[uint32]*Session)
	r.createSessionChan = make(chan *Session, 1000)
	r.destroySessionChan = make(chan *Session, 1000)
	r.globalGsOnlineMap = make(map[uint32]string)
	r.kcpEventInput = make(chan *KcpEvent, 1000)
	r.kcpEventOutput = make(chan *KcpEvent, 1000)
	r.reLoginRemoteKickRegChan = make(chan *RemoteKick, 1000)
	r.serverCmdProtoMap = cmd.NewCmdProtoMap()
	if config.GetConfig().Hk4e.ClientProtoProxyEnable {
		r.clientCmdProtoMap = client_proto.NewClientCmdProtoMap()
	}
	r.messageQueue = messageQueue
	r.run()
	return r
}

func (k *KcpConnectManager) run() {
	// 读取密钥相关文件
	k.signRsaKey, k.encRsaKeyMap, _ = region.LoadRsaKey()
	// key
	rsp, err := k.discovery.GetRegionEc2B(context.TODO(), &api.NullMsg{})
	if err != nil {
		logger.Error("get region ec2b error: %v", err)
		return
	}
	ec2b, err := random.LoadEc2bKey(rsp.Data)
	if err != nil {
		logger.Error("parse region ec2b error: %v", err)
		return
	}
	regionEc2b := random.NewEc2b()
	regionEc2b.SetSeed(ec2b.Seed())
	// 3.7的时候修改了xor 首包无需使用加密密钥
	if k.getGateMaxVersion() < 370 {
		k.dispatchKey = regionEc2b.XorKey()
	} else {
		// 全部填充为0 不然会出问题
		k.dispatchKey = make([]byte, 4096)
	}
	// kcp
	port := strconv.Itoa(int(config.GetConfig().Hk4e.KcpPort))
	listener, err := kcp.ListenWithOptions("0.0.0.0:" + port)
	if err != nil {
		logger.Error("listen kcp err: %v", err)
		return
	}
	go k.enetHandle(listener)
	go k.eventHandle()
	go k.sendMsgHandle()
	go k.acceptHandle(listener)
	go k.gateNetInfo()
	k.syncGlobalGsOnlineMap()
	go k.autoSyncGlobalGsOnlineMap()
}

func (k *KcpConnectManager) Close() {
	k.closeAllKcpConn()
}

// getGateMaxVersion 获取gate最大可兼容的版本
func (k *KcpConnectManager) getGateMaxVersion() (maxVersion int) {
	versionSplit := strings.Split(config.GetConfig().Hk4e.Version, ",")
	for _, verStr := range versionSplit {
		version, err := strconv.Atoi(verStr)
		if err != nil {
			logger.Error("version a to i error: %v", err)
			return
		}
		if version > maxVersion {
			maxVersion = version
		}
	}
	return
}

func (k *KcpConnectManager) gateNetInfo() {
	ticker := time.NewTicker(time.Second * 60)
	kcpErrorCount := uint64(0)
	for {
		<-ticker.C
		snmp := kcp.DefaultSnmp.Copy()
		kcpErrorCount += snmp.KCPInErrors
		logger.Info("kcp send: %v B/s, kcp recv: %v B/s", snmp.BytesSent/60, snmp.BytesReceived/60)
		logger.Info("udp send: %v B/s, udp recv: %v B/s", snmp.OutBytes/60, snmp.InBytes/60)
		logger.Info("udp send: %v pps, udp recv: %v pps", snmp.OutPkts/60, snmp.InPkts/60)
		clientConnNum := atomic.LoadInt32(&CLIENT_CONN_NUM)
		logger.Info("conn num: %v, new conn num: %v, kcp error num: %v", clientConnNum, snmp.CurrEstab, kcpErrorCount)
		kcp.DefaultSnmp.Reset()
	}
}

// 接收并创建新连接处理函数
func (k *KcpConnectManager) acceptHandle(listener *kcp.Listener) {
	logger.Info("accept handle start")
	for {
		conn, err := listener.AcceptKCP()
		if err != nil {
			logger.Error("accept kcp err: %v", err)
			return
		}
		convId := conn.GetConv()
		if k.openState == false {
			logger.Error("gate not open, convId: %v", convId)
			_ = conn.Close()
			continue
		}
		conn.SetACKNoDelay(true)
		conn.SetWriteDelay(false)
		conn.SetWindowSize(255, 255)
		atomic.AddInt32(&CLIENT_CONN_NUM, 1)
		logger.Info("client connect, convId: %v", convId)
		kcpRawSendChan := make(chan *ProtoMsg, 1000)
		session := &Session{
			conn:                   conn,
			connState:              ConnEst,
			userId:                 0,
			kcpRawSendChan:         kcpRawSendChan,
			seed:                   0,
			xorKey:                 k.dispatchKey,
			changeXorKeyFin:        false,
			gsServerAppId:          "",
			anticheatServerAppId:   "",
			pathfindingServerAppId: "",
			useMagicSeed:           false,
		}
		go k.recvHandle(session)
		go k.sendHandle(session)
		// 连接建立成功通知
		k.kcpEventOutput <- &KcpEvent{
			ConvId:       convId,
			EventId:      KcpConnEstNotify,
			EventMessage: conn.RemoteAddr().String(),
		}
	}
}

// 连接事件处理函数
func (k *KcpConnectManager) enetHandle(listener *kcp.Listener) {
	logger.Info("enet handle start")
	// conv短时间内唯一生成
	convGenMap := make(map[uint64]int64)
	pktFreqLimitCounter := 0
	pktFreqLimitTimer := time.Now().UnixNano()
	for {
		enetNotify := <-listener.EnetNotify
		logger.Info("[Enet Notify], addr: %v, conv: %v, conn: %v, enet: %v", enetNotify.Addr, enetNotify.ConvId, enetNotify.ConnType, enetNotify.EnetType)
		switch enetNotify.ConnType {
		case kcp.ConnEnetSyn:
			// 连接建立握手包频率限制
			pktFreqLimitCounter++
			if pktFreqLimitCounter > ConnSynPacketFreqLimit {
				now := time.Now().UnixNano()
				if now-pktFreqLimitTimer > int64(time.Second) {
					pktFreqLimitCounter = 0
					pktFreqLimitTimer = now
				} else {
					continue
				}
			}
			if enetNotify.EnetType != kcp.EnetClientConnectKey {
				continue
			}
			// 清理老旧的conv
			now := time.Now().UnixNano()
			oldConvList := make([]uint64, 0)
			for conv, timestamp := range convGenMap {
				if now-timestamp > int64(time.Hour) {
					oldConvList = append(oldConvList, conv)
				}
			}
			delConvList := make([]uint64, 0)
			k.sessionMapLock.RLock()
			for _, conv := range oldConvList {
				_, exist := k.sessionConvIdMap[conv]
				if !exist {
					delConvList = append(delConvList, conv)
					delete(convGenMap, conv)
				}
			}
			k.sessionMapLock.RUnlock()
			logger.Info("clean dead conv list: %v", delConvList)
			// 生成没用过的conv
			var conv uint64
			for {
				convData := random.GetRandomByte(8)
				convDataBuffer := bytes.NewBuffer(convData)
				_ = binary.Read(convDataBuffer, binary.LittleEndian, &conv)
				_, exist := convGenMap[conv]
				if exist {
					continue
				} else {
					convGenMap[conv] = time.Now().UnixNano()
					break
				}
			}
			listener.SendEnetNotifyToPeer(&kcp.Enet{
				Addr:     enetNotify.Addr,
				ConvId:   conv,
				ConnType: kcp.ConnEnetEst,
				EnetType: enetNotify.EnetType,
			})
		case kcp.ConnEnetEst:
		case kcp.ConnEnetFin:
			session := k.GetSessionByConvId(enetNotify.ConvId)
			if session == nil {
				logger.Error("session not exist, conv: %v", enetNotify.ConvId)
				continue
			}
			session.conn.SendEnetNotifyToPeer(&kcp.Enet{
				ConnType: kcp.ConnEnetFin,
				EnetType: enetNotify.EnetType,
			})
			_ = session.conn.Close()
		case kcp.ConnEnetAddrChange:
			// 连接地址改变通知
			k.kcpEventOutput <- &KcpEvent{
				ConvId:       enetNotify.ConvId,
				EventId:      KcpConnAddrChangeNotify,
				EventMessage: enetNotify.Addr,
			}
		default:
		}
	}
}

// Session 连接会话结构 只允许定义并发安全或者简单的基础数据结构
type Session struct {
	conn                   *kcp.UDPSession
	connState              uint8
	userId                 uint32
	kcpRawSendChan         chan *ProtoMsg
	seed                   uint64
	xorKey                 []byte
	changeXorKeyFin        bool
	gsServerAppId          string
	anticheatServerAppId   string
	pathfindingServerAppId string
	useMagicSeed           bool
}

// 接收
func (k *KcpConnectManager) recvHandle(session *Session) {
	logger.Info("recv handle start")
	conn := session.conn
	convId := conn.GetConv()
	recvBuf := make([]byte, PacketMaxLen)
	pktFreqLimitCounter := 0
	pktFreqLimitTimer := time.Now().UnixNano()
	for {
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * ConnRecvTimeout))
		recvLen, err := conn.Read(recvBuf)
		if err != nil {
			logger.Error("exit recv loop, conn read err: %v, convId: %v", err, convId)
			k.closeKcpConn(session, kcp.EnetServerKick)
			break
		}
		// 收包频率限制
		pktFreqLimitCounter++
		if pktFreqLimitCounter > RecvPacketFreqLimit {
			now := time.Now().UnixNano()
			if now-pktFreqLimitTimer > int64(time.Second) {
				pktFreqLimitCounter = 0
				pktFreqLimitTimer = now
			} else {
				logger.Error("exit recv loop, client packet send freq too high, convId: %v, pps: %v", convId, pktFreqLimitCounter)
				k.closeKcpConn(session, kcp.EnetPacketFreqTooHigh)
				break
			}
		}
		recvData := recvBuf[:recvLen]
		kcpMsgList := make([]*KcpMsg, 0)
		DecodeBinToPayload(recvData, convId, &kcpMsgList, session.xorKey)
		for _, v := range kcpMsgList {
			protoMsgList := ProtoDecode(v, k.serverCmdProtoMap, k.clientCmdProtoMap)
			for _, vv := range protoMsgList {
				k.recvMsgHandle(vv, session)
			}
		}
	}
}

// 发送
func (k *KcpConnectManager) sendHandle(session *Session) {
	logger.Info("send handle start")
	conn := session.conn
	convId := conn.GetConv()
	pktFreqLimitCounter := 0
	pktFreqLimitTimer := time.Now().UnixNano()
	for {
		protoMsg, ok := <-session.kcpRawSendChan
		if !ok {
			logger.Error("exit send loop, send chan close, convId: %v", convId)
			k.closeKcpConn(session, kcp.EnetServerKick)
			break
		}
		kcpMsg := ProtoEncode(protoMsg, k.serverCmdProtoMap, k.clientCmdProtoMap)
		if kcpMsg == nil {
			logger.Error("decode kcp msg is nil, convId: %v", convId)
			continue
		}
		bin := EncodePayloadToBin(kcpMsg, session.xorKey)
		_ = conn.SetWriteDeadline(time.Now().Add(time.Second * ConnSendTimeout))
		_, err := conn.Write(bin)
		if err != nil {
			logger.Error("exit send loop, conn write err: %v, convId: %v", err, convId)
			k.closeKcpConn(session, kcp.EnetServerKick)
			break
		}
		// 发包频率限制
		pktFreqLimitCounter++
		if pktFreqLimitCounter > SendPacketFreqLimit {
			now := time.Now().UnixNano()
			if now-pktFreqLimitTimer > int64(time.Second) {
				pktFreqLimitCounter = 0
				pktFreqLimitTimer = now
			} else {
				logger.Error("exit send loop, server packet send freq too high, convId: %v, pps: %v", convId, pktFreqLimitCounter)
				k.closeKcpConn(session, kcp.EnetPacketFreqTooHigh)
				break
			}
		}
		if session.changeXorKeyFin == false && protoMsg.CmdId == cmd.GetPlayerTokenRsp {
			// XOR密钥切换
			logger.Info("change session xor key, convId: %v", convId)
			session.changeXorKeyFin = true
			keyBlock := random.NewKeyBlock(session.seed, session.useMagicSeed)
			xorKey := keyBlock.XorKey()
			key := make([]byte, 4096)
			copy(key, xorKey[:])
			session.xorKey = key
		}
	}
}

// 强制关闭指定连接
func (k *KcpConnectManager) forceCloseKcpConn(convId uint64, reason uint32) {
	session := k.GetSessionByConvId(convId)
	if session == nil {
		logger.Error("session not exist, convId: %v", convId)
		return
	}
	k.closeKcpConn(session, reason)
	logger.Info("conn has been force close, convId: %v", convId)
}

// 关闭指定连接
func (k *KcpConnectManager) closeKcpConn(session *Session, enetType uint32) {
	if session.connState == ConnClose {
		return
	}
	session.connState = ConnClose
	conn := session.conn
	convId := conn.GetConv()
	// 清理数据
	k.DeleteSession(session.conn.GetConv(), session.userId)
	// 关闭连接
	err := conn.Close()
	if err == nil {
		conn.SendEnetNotifyToPeer(&kcp.Enet{
			ConnType: kcp.ConnEnetFin,
			EnetType: enetType,
		})
	}
	// 连接关闭通知
	k.kcpEventOutput <- &KcpEvent{
		ConvId:  convId,
		EventId: KcpConnCloseNotify,
	}
	// 通知GS玩家下线
	connCtrlMsg := new(mq.ConnCtrlMsg)
	connCtrlMsg.UserId = session.userId
	k.messageQueue.SendToGs(session.gsServerAppId, &mq.NetMsg{
		MsgType:     mq.MsgTypeConnCtrl,
		EventId:     mq.UserOfflineNotify,
		ConnCtrlMsg: connCtrlMsg,
	})
	logger.Info("send to gs user offline, ConvId: %v, UserId: %v", convId, connCtrlMsg.UserId)
	k.destroySessionChan <- session
	atomic.AddInt32(&CLIENT_CONN_NUM, -1)
}

// 关闭所有连接
func (k *KcpConnectManager) closeAllKcpConn() {
	sessionList := make([]*Session, 0)
	k.sessionMapLock.RLock()
	for _, session := range k.sessionConvIdMap {
		sessionList = append(sessionList, session)
	}
	k.sessionMapLock.RUnlock()
	for _, session := range sessionList {
		k.closeKcpConn(session, kcp.EnetServerShutdown)
	}
}

func (k *KcpConnectManager) GetSessionByConvId(convId uint64) *Session {
	k.sessionMapLock.RLock()
	session, _ := k.sessionConvIdMap[convId]
	k.sessionMapLock.RUnlock()
	return session
}

func (k *KcpConnectManager) GetSessionByUserId(userId uint32) *Session {
	k.sessionMapLock.RLock()
	session, _ := k.sessionUserIdMap[userId]
	k.sessionMapLock.RUnlock()
	return session
}

func (k *KcpConnectManager) SetSession(session *Session, convId uint64, userId uint32) {
	k.sessionMapLock.Lock()
	k.sessionConvIdMap[convId] = session
	k.sessionUserIdMap[userId] = session
	k.sessionMapLock.Unlock()
}

func (k *KcpConnectManager) DeleteSession(convId uint64, userId uint32) {
	k.sessionMapLock.Lock()
	delete(k.sessionConvIdMap, convId)
	delete(k.sessionUserIdMap, userId)
	k.sessionMapLock.Unlock()
}

func (k *KcpConnectManager) autoSyncGlobalGsOnlineMap() {
	ticker := time.NewTicker(time.Second * 60)
	for {
		<-ticker.C
		k.syncGlobalGsOnlineMap()
	}
}

func (k *KcpConnectManager) syncGlobalGsOnlineMap() {
	rsp, err := k.discovery.GetGlobalGsOnlineMap(context.TODO(), nil)
	if err != nil {
		logger.Error("get global gs online map error: %v", err)
		return
	}
	copyMap := make(map[uint32]string)
	for k, v := range rsp.GlobalGsOnlineMap {
		copyMap[k] = v
	}
	copyMapLen := len(copyMap)
	k.globalGsOnlineMapLock.Lock()
	k.globalGsOnlineMap = copyMap
	k.globalGsOnlineMapLock.Unlock()
	logger.Info("sync global gs online map finish, len: %v", copyMapLen)
}

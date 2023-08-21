package game

import (
	"encoding/base64"

	"hk4e/gdconf"
	"hk4e/gs/model"
	"hk4e/pkg/logger"
	"hk4e/pkg/random"
	"hk4e/protocol/cmd"
	"hk4e/protocol/proto"

	"google.golang.org/protobuf/encoding/protojson"
	pb "google.golang.org/protobuf/proto"
)

// GM函数模块
// GM函数只支持基本类型的简单参数传入

type GMCmd struct {
}

// 玩家通用GM指令

// GMTeleportPlayer 传送玩家
func (g *GMCmd) GMTeleportPlayer(userId, sceneId uint32, posX, posY, posZ float64) {
	player := USER_MANAGER.GetOnlineUser(userId)
	if player == nil {
		logger.Error("player is nil, uid: %v", userId)
		return
	}
	GAME.TeleportPlayer(
		player,
		proto.EnterReason_ENTER_REASON_GM,
		sceneId,
		&model.Vector{X: posX, Y: posY, Z: posZ},
		new(model.Vector),
		0,
		0,
	)
}

// GMAddUserItem 给予玩家物品
func (g *GMCmd) GMAddUserItem(userId, itemId, itemCount uint32) {
	GAME.AddUserItem(userId, []*ChangeItem{
		{
			ItemId:      itemId,
			ChangeCount: itemCount,
		},
	}, true, 0)
}

// GMAddUserWeapon 给予玩家武器
func (g *GMCmd) GMAddUserWeapon(userId, itemId, itemCount uint32) {
	// 武器数量
	for i := uint32(0); i < itemCount; i++ {
		// 给予武器
		GAME.AddUserWeapon(userId, itemId)
	}
}

// GMAddUserReliquary 给予玩家圣遗物
func (g *GMCmd) GMAddUserReliquary(userId, itemId, itemCount uint32) {
	// 圣遗物数量
	for i := uint32(0); i < itemCount; i++ {
		// 给予圣遗物
		GAME.AddUserReliquary(userId, itemId)
	}
}

// GMAddUserAvatar 给予玩家角色
func (g *GMCmd) GMAddUserAvatar(userId, avatarId uint32) {
	// 添加角色
	GAME.AddUserAvatar(userId, avatarId)
	// TODO 设置角色 等以后做到角色升级之类的再说
	// avatar := player.AvatarMap[avatarId]
}

// GMAddUserCostume 给予玩家时装
func (g *GMCmd) GMAddUserCostume(userId, costumeId uint32) {
	// 添加时装
	GAME.AddUserCostume(userId, costumeId)
}

// GMAddUserFlycloak 给予玩家风之翼
func (g *GMCmd) GMAddUserFlycloak(userId, flycloakId uint32) {
	// 添加风之翼
	GAME.AddUserFlycloak(userId, flycloakId)
}

// GMAddUserAllItem 给予玩家所有物品
func (g *GMCmd) GMAddUserAllItem(userId, itemCount uint32) {
	GAME.LogoutPlayer(userId)
	itemList := make([]*ChangeItem, 0)
	for itemId := range GAME.GetAllItemDataConfig() {
		itemList = append(itemList, &ChangeItem{
			ItemId:      uint32(itemId),
			ChangeCount: itemCount,
		})
	}
	GAME.AddUserItem(userId, itemList, false, 0)
}

// GMAddUserAllWeapon 给予玩家所有武器
func (g *GMCmd) GMAddUserAllWeapon(userId, itemCount uint32) {
	for itemId := range GAME.GetAllWeaponDataConfig() {
		g.GMAddUserWeapon(userId, uint32(itemId), itemCount)
	}
}

// GMAddUserAllReliquary 给予玩家所有圣遗物
func (g *GMCmd) GMAddUserAllReliquary(userId, itemCount uint32) {
	GAME.LogoutPlayer(userId)
	for itemId := range GAME.GetAllReliquaryDataConfig() {
		g.GMAddUserReliquary(userId, uint32(itemId), itemCount)
	}
}

// GMAddUserAllAvatar 给予玩家所有角色
func (g *GMCmd) GMAddUserAllAvatar(userId uint32) {
	for avatarId := range GAME.GetAllAvatarDataConfig() {
		g.GMAddUserAvatar(userId, uint32(avatarId))
	}
}

// GMAddUserAllCostume 给予玩家所有时装
func (g *GMCmd) GMAddUserAllCostume(userId uint32) {
	for costumeId := range gdconf.GetAvatarCostumeDataMap() {
		g.GMAddUserCostume(userId, uint32(costumeId))
	}
}

// GMAddUserAllFlycloak 给予玩家所有风之翼
func (g *GMCmd) GMAddUserAllFlycloak(userId uint32) {
	for flycloakId := range gdconf.GetAvatarFlycloakDataMap() {
		g.GMAddUserFlycloak(userId, uint32(flycloakId))
	}
}

// GMAddUserAllEvery 给予玩家所有内容
func (g *GMCmd) GMAddUserAllEvery(userId, itemCount uint32) {
	GAME.LogoutPlayer(userId)
	// 给予玩家所有物品
	g.GMAddUserAllItem(userId, itemCount)
	// 给予玩家所有武器
	g.GMAddUserAllWeapon(userId, itemCount)
	// 给予玩家所有圣遗物
	g.GMAddUserAllReliquary(userId, itemCount)
	// 给予玩家所有角色
	g.GMAddUserAllAvatar(userId)
	// 给予玩家所有时装
	g.GMAddUserAllCostume(userId)
	// 给予玩家所有风之翼
	g.GMAddUserAllFlycloak(userId)
}

// GMAddQuest 添加任务
func (g *GMCmd) GMAddQuest(userId uint32, questId uint32) {
	player := USER_MANAGER.GetOnlineUser(userId)
	if player == nil {
		logger.Error("player is nil, uid: %v", userId)
		return
	}
	dbQuest := player.GetDbQuest()
	dbQuest.AddQuest(questId)
	dbQuest.StartQuest(questId)
	ntf := &proto.QuestListUpdateNotify{
		QuestList: make([]*proto.Quest, 0),
	}
	ntf.QuestList = append(ntf.QuestList, GAME.PacketQuest(player, questId))
	GAME.SendMsg(cmd.QuestListUpdateNotify, player.PlayerID, player.ClientSeq, ntf)
}

// GMFinishQuest 完成任务
func (g *GMCmd) GMFinishQuest(userId uint32, questId uint32) {
	player := USER_MANAGER.GetOnlineUser(userId)
	if player == nil {
		logger.Error("player is nil, uid: %v", userId)
		return
	}
	dbQuest := player.GetDbQuest()
	dbQuest.ForceFinishQuest(questId)
	ntf := &proto.QuestListUpdateNotify{
		QuestList: make([]*proto.Quest, 0),
	}
	ntf.QuestList = append(ntf.QuestList, GAME.PacketQuest(player, questId))
	GAME.SendMsg(cmd.QuestListUpdateNotify, player.PlayerID, player.ClientSeq, ntf)
	GAME.AcceptQuest(player, true)
}

// GMForceFinishAllQuest 强制完成当前所有任务
func (g *GMCmd) GMForceFinishAllQuest(userId uint32) {
	player := USER_MANAGER.GetOnlineUser(userId)
	if player == nil {
		logger.Error("player is nil, uid: %v", userId)
		return
	}
	dbQuest := player.GetDbQuest()
	ntf := &proto.QuestListUpdateNotify{
		QuestList: make([]*proto.Quest, 0),
	}
	for _, quest := range dbQuest.GetQuestMap() {
		dbQuest.ForceFinishQuest(quest.QuestId)
		pbQuest := GAME.PacketQuest(player, quest.QuestId)
		if pbQuest == nil {
			continue
		}
		ntf.QuestList = append(ntf.QuestList, pbQuest)
	}
	GAME.SendMsg(cmd.QuestListUpdateNotify, player.PlayerID, player.ClientSeq, ntf)
	GAME.AcceptQuest(player, true)
}

// GMUnlockAllPoint 解锁场景全部传送点
func (g *GMCmd) GMUnlockAllPoint(userId uint32, sceneId uint32) {
	player := USER_MANAGER.GetOnlineUser(userId)
	if player == nil {
		logger.Error("player is nil, uid: %v", userId)
		return
	}
	dbWorld := player.GetDbWorld()
	dbScene := dbWorld.GetSceneById(sceneId)
	if dbScene == nil {
		logger.Error("db scene is nil, uid: %v", sceneId)
		return
	}
	scenePointMapConfig := gdconf.GetScenePointMapBySceneId(int32(sceneId))
	for _, pointData := range scenePointMapConfig {
		dbScene.UnlockPoint(uint32(pointData.Id))
	}
	GAME.SendMsg(cmd.ScenePointUnlockNotify, player.PlayerID, player.ClientSeq, &proto.ScenePointUnlockNotify{
		SceneId:         sceneId,
		PointList:       dbScene.GetUnlockPointList(),
		UnhidePointList: nil,
	})
}

// GMCreateMonster 在玩家附近创建怪物
func (g *GMCmd) GMCreateMonster(userId uint32, monsterId uint32) {
	player := USER_MANAGER.GetOnlineUser(userId)
	if player == nil {
		logger.Error("player is nil, uid: %v", userId)
		return
	}
	GAME.CreateMonster(player, nil, monsterId)
}

// GMCreateGadget 在玩家附近创建物件
func (g *GMCmd) GMCreateGadget(userId uint32, gadgetId uint32) {
	player := USER_MANAGER.GetOnlineUser(userId)
	if player == nil {
		logger.Error("player is nil, uid: %v", userId)
		return
	}
	GAME.CreateGadget(player, nil, gadgetId, nil)
}

// 系统级GM指令

func (g *GMCmd) ChangePlayerCmdPerm(userId uint32, cmdPerm uint8) {
	player := USER_MANAGER.GetOnlineUser(userId)
	if player == nil {
		logger.Error("player is nil, uid: %v", userId)
		return
	}
	player.CmdPerm = cmdPerm
}

func (g *GMCmd) ReloadGameDataConfig() {
	LOCAL_EVENT_MANAGER.GetLocalEventChan() <- &LocalEvent{
		EventId: ReloadGameDataConfig,
	}
}

func (g *GMCmd) XLuaDebug(userId uint32, luacBase64 string) {
	logger.Debug("xlua debug, uid: %v, luac: %v", userId, luacBase64)
	player := USER_MANAGER.GetOnlineUser(userId)
	if player == nil {
		logger.Error("player is nil, uid: %v", userId)
		return
	}
	// 只有在线玩家主动开启之后才能发送
	if !player.XLuaDebug {
		logger.Error("player xlua debug not enable, uid: %v", userId)
		return
	}
	luac, err := base64.StdEncoding.DecodeString(luacBase64)
	if err != nil {
		logger.Error("decode luac error: %v", err)
		return
	}
	GAME.SendMsg(cmd.WindSeedClientNotify, player.PlayerID, 0, &proto.WindSeedClientNotify{
		Notify: &proto.WindSeedClientNotify_AreaNotify_{
			AreaNotify: &proto.WindSeedClientNotify_AreaNotify{
				AreaCode: luac,
				AreaId:   1,
				AreaType: 1,
			},
		},
	})
}

func (g *GMCmd) PlayAudio() {
	PlayAudio()
}

func (g *GMCmd) UpdateFrame(rgb bool) {
	UpdateFrame(rgb)
}

var RobotUidCounter uint32 = 0

func (g *GMCmd) CreateRobotInBigWorld(uid uint32, name string, avatarId uint32) {
	if !GAME.IsMainGs() {
		return
	}
	if uid == 0 {
		RobotUidCounter++
		uid = 1000000 + RobotUidCounter
	}
	if name == "" {
		name = random.GetRandomStr(8)
	}
	if avatarId == 0 {
		for _, avatarData := range gdconf.GetAvatarDataMap() {
			avatarId = uint32(avatarData.AvatarId)
			break
		}
	}
	aiWorld := WORLD_MANAGER.GetAiWorld()
	robot := GAME.CreateRobot(uid, name, name)
	GAME.AddUserAvatar(uid, avatarId)
	dbAvatar := robot.GetDbAvatar()
	GAME.SetUpAvatarTeamReq(robot, &proto.SetUpAvatarTeamReq{
		TeamId:             1,
		AvatarTeamGuidList: []uint64{dbAvatar.AvatarMap[avatarId].Guid},
		CurAvatarGuid:      dbAvatar.AvatarMap[avatarId].Guid,
	})
	GAME.SetPlayerHeadImageReq(robot, &proto.SetPlayerHeadImageReq{
		AvatarId: avatarId,
	})
	GAME.JoinPlayerSceneReq(robot, &proto.JoinPlayerSceneReq{
		TargetUid: aiWorld.owner.PlayerID,
	})
	GAME.EnterSceneReadyReq(robot, &proto.EnterSceneReadyReq{
		EnterSceneToken: aiWorld.GetEnterSceneToken(),
	})
	GAME.SceneInitFinishReq(robot, &proto.SceneInitFinishReq{
		EnterSceneToken: aiWorld.GetEnterSceneToken(),
	})
	GAME.EnterSceneDoneReq(robot, &proto.EnterSceneDoneReq{
		EnterSceneToken: aiWorld.GetEnterSceneToken(),
	})
	GAME.PostEnterSceneReq(robot, &proto.PostEnterSceneReq{
		EnterSceneToken: aiWorld.GetEnterSceneToken(),
	})
	activeAvatarId := aiWorld.GetPlayerActiveAvatarId(robot)
	pos := new(model.Vector)
	rot := new(model.Vector)
	for _, targetPlayer := range aiWorld.GetAllPlayer() {
		if targetPlayer.PlayerID < PlayerBaseUid {
			continue
		}
		pos = &model.Vector{X: targetPlayer.Pos.X, Y: targetPlayer.Pos.Y, Z: targetPlayer.Pos.Z}
		rot = &model.Vector{X: targetPlayer.Rot.X, Y: targetPlayer.Rot.Y, Z: targetPlayer.Rot.Z}
	}
	entityMoveInfo := &proto.EntityMoveInfo{
		EntityId: aiWorld.GetPlayerWorldAvatarEntityId(robot, activeAvatarId),
		MotionInfo: &proto.MotionInfo{
			Pos:   &proto.Vector{X: float32(pos.X), Y: float32(pos.Y), Z: float32(pos.Z)},
			Rot:   &proto.Vector{X: float32(rot.X), Y: float32(rot.Y), Z: float32(rot.Z)},
			State: proto.MotionState_MOTION_STANDBY,
		},
		SceneTime:   0,
		ReliableSeq: 0,
	}
	combatData, err := pb.Marshal(entityMoveInfo)
	if err != nil {
		return
	}
	GAME.CombatInvocationsNotify(robot, &proto.CombatInvocationsNotify{
		InvokeList: []*proto.CombatInvokeEntry{{
			CombatData:   combatData,
			ForwardType:  proto.ForwardType_FORWARD_TO_ALL_EXCEPT_CUR,
			ArgumentType: proto.CombatTypeArgument_ENTITY_MOVE,
		}},
	})
	GAME.UnionCmdNotify(robot, &proto.UnionCmdNotify{})
}

func (g *GMCmd) ServerAnnounce(announceId uint32, announceMsg string, isRevoke bool) {
	if !isRevoke {
		GAME.ServerAnnounceNotify(announceId, announceMsg)
	} else {
		GAME.ServerAnnounceRevokeNotify(announceId)
	}
}

func (g *GMCmd) SendMsgToPlayer(cmdName string, userId uint32, msgJson string) {
	if cmdProtoMap == nil {
		cmdProtoMap = cmd.NewCmdProtoMap()
	}
	cmdId := cmdProtoMap.GetCmdIdByCmdName(cmdName)
	if cmdId == 0 {
		logger.Error("cmd name not found")
		return
	}
	if cmdId == cmd.WindSeedClientNotify {
		logger.Error("what are you doing ???")
		return
	}
	msg := cmdProtoMap.GetProtoObjByCmdId(cmdId)
	err := protojson.Unmarshal([]byte(msgJson), msg)
	if err != nil {
		logger.Error("parse msg error: %v", err)
		return
	}
	GAME.SendMsg(cmdId, userId, 0, msg)
}

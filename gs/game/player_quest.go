package game

import (
	"strconv"
	"strings"

	"hk4e/common/constant"
	"hk4e/gdconf"
	"hk4e/gs/model"
	"hk4e/pkg/logger"
	"hk4e/protocol/cmd"
	"hk4e/protocol/proto"

	pb "google.golang.org/protobuf/proto"
)

// AddQuestContentProgressReq 添加任务内容进度请求
func (g *Game) AddQuestContentProgressReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.AddQuestContentProgressReq)
	logger.Debug("AddQuestContentProgressReq: %v", req)

	g.TriggerQuest(player, int32(req.ContentType), "", int32(req.Param))

	// g.AddQuestProgress(player, req)

	rsp := &proto.AddQuestContentProgressRsp{
		ContentType: req.ContentType,
	}
	g.SendMsg(cmd.AddQuestContentProgressRsp, player.PlayerID, player.ClientSeq, rsp)

	g.AcceptQuest(player, true)
}

// AddQuestProgress 添加任务进度
func (g *Game) AddQuestProgress(player *model.Player, req *proto.AddQuestContentProgressReq) {
	dbQuest := player.GetDbQuest()
	updateQuestIdList := make([]uint32, 0)
	for _, quest := range dbQuest.GetQuestMap() {
		questDataConfig := gdconf.GetQuestDataById(int32(quest.QuestId))
		if questDataConfig == nil {
			logger.Error("get quest data config is nil, questId: %v", quest.QuestId)
			continue
		}
		for index, finishCond := range questDataConfig.FinishCondList {
			if len(finishCond.Param) != 1 {
				continue
			}
			if req.ContentType != uint32(finishCond.Type) || req.Param != uint32(finishCond.Param[0]) {
				continue
			}
			dbQuest.AddQuestProgress(quest.QuestId, index, req.AddProgress)
			updateQuestIdList = append(updateQuestIdList, quest.QuestId)
		}
	}
	for _, questId := range updateQuestIdList {
		quest := dbQuest.GetQuestById(questId)
		if quest == nil {
			logger.Error("get quest is nil, questId: %v", quest.QuestId)
			continue
		}
		ntf := &proto.QuestProgressUpdateNotify{
			QuestId:            quest.QuestId,
			FinishProgressList: quest.FinishProgressList,
		}
		g.SendMsg(cmd.QuestProgressUpdateNotify, player.PlayerID, player.ClientSeq, ntf)
	}
}

// AcceptQuest 接取当前条件下能接取到的全部任务
func (g *Game) AcceptQuest(player *model.Player, notifyClient bool) {
	dbQuest := player.GetDbQuest()
	addQuestIdList := make([]uint32, 0)
	for _, questData := range gdconf.GetQuestDataMap() {
		if dbQuest.GetQuestById(uint32(questData.QuestId)) != nil {
			continue
		}
		acceptCondResultList := make([]bool, 0)
		for _, acceptCond := range questData.AcceptCondList {
			result := false
			switch acceptCond.Type {
			case constant.QUEST_ACCEPT_COND_TYPE_STATE_EQUAL:
				// 某个任务状态等于 参数1:任务id 参数2:任务状态
				if len(acceptCond.Param) != 2 {
					break
				}
				quest := dbQuest.GetQuestById(uint32(acceptCond.Param[0]))
				if quest == nil {
					break
				}
				if quest.State != uint8(acceptCond.Param[1]) {
					break
				}
				result = true
			case constant.QUEST_ACCEPT_COND_TYPE_STATE_NOT_EQUAL:
				// 某个任务状态不等于 参数1:任务id 参数2:任务状态
				if len(acceptCond.Param) != 2 {
					break
				}
				quest := dbQuest.GetQuestById(uint32(acceptCond.Param[0]))
				if quest == nil {
					break
				}
				if quest.State == uint8(acceptCond.Param[1]) {
					break
				}
				result = true
			default:
				break
			}
			acceptCondResultList = append(acceptCondResultList, result)
		}
		canAccept := false
		switch questData.AcceptCondCompose {
		case constant.QUEST_LOGIC_TYPE_NONE:
			fallthrough
		case constant.QUEST_LOGIC_TYPE_AND:
			canAccept = true
			for _, acceptCondResult := range acceptCondResultList {
				if !acceptCondResult {
					canAccept = false
					break
				}
			}
		case constant.QUEST_LOGIC_TYPE_OR:
			canAccept = false
			for _, acceptCondResult := range acceptCondResultList {
				if acceptCondResult {
					canAccept = true
					break
				}
			}
		}
		if canAccept {
			dbQuest.AddQuest(uint32(questData.QuestId))
			addQuestIdList = append(addQuestIdList, uint32(questData.QuestId))
		}
	}
	if notifyClient {
		ntf := &proto.QuestListUpdateNotify{
			QuestList: make([]*proto.Quest, 0),
		}
		for _, questId := range addQuestIdList {
			pbQuest := g.PacketQuest(player, questId)
			if pbQuest == nil {
				continue
			}
			ntf.QuestList = append(ntf.QuestList, pbQuest)
		}
		g.SendMsg(cmd.QuestListUpdateNotify, player.PlayerID, player.ClientSeq, ntf)
	}
	// TODO 判断任务是否能开始
	for _, questId := range addQuestIdList {
		g.StartQuest(player, questId, notifyClient)
	}
}

// StartQuest 开始一个任务
func (g *Game) StartQuest(player *model.Player, questId uint32, notifyClient bool) {
	dbQuest := player.GetDbQuest()
	dbQuest.StartQuest(questId)

	g.QuestExec(player, questId)
	g.QuestStartTriggerCheck(player, questId)

	if notifyClient {
		ntf := &proto.QuestListUpdateNotify{
			QuestList: make([]*proto.Quest, 0),
		}
		pbQuest := g.PacketQuest(player, questId)
		if pbQuest == nil {
			return
		}
		ntf.QuestList = append(ntf.QuestList, pbQuest)
		g.SendMsg(cmd.QuestListUpdateNotify, player.PlayerID, player.ClientSeq, ntf)
	}
}

// QuestExec 任务开始执行触发操作
func (g *Game) QuestExec(player *model.Player, questId uint32) {
	questDataConfig := gdconf.GetQuestDataById(int32(questId))
	if questDataConfig == nil {
		return
	}
	for _, questExec := range questDataConfig.StartExecList {
		switch questExec.Type {
		case constant.QUEST_EXEC_TYPE_NOTIFY_GROUP_LUA:
		case constant.QUEST_EXEC_TYPE_REFRESH_GROUP_SUITE:
			if len(questExec.Param) != 2 {
				continue
			}
			split := strings.Split(questExec.Param[1], ",")
			if len(split) != 2 {
				continue
			}
			groupId, err := strconv.Atoi(split[0])
			if err != nil {
				continue
			}
			suiteId, err := strconv.Atoi(split[1])
			if err != nil {
				continue
			}
			g.AddSceneGroupSuite(player, uint32(groupId), uint8(suiteId))
		}
	}
}

// 通用参数匹配
func matchParamEqual(param1 []int32, param2 []int32, num int) bool {
	if len(param1) != num || len(param2) != num {
		return false
	}
	for i := 0; i < num; i++ {
		if param1[i] != param2[i] {
			return false
		}
	}
	return true
}

// TriggerQuest 触发任务
func (g *Game) TriggerQuest(player *model.Player, cond int32, complexParam string, param ...int32) {
	dbQuest := player.GetDbQuest()
	updateQuestIdList := make([]uint32, 0)
	for _, quest := range dbQuest.GetQuestMap() {
		questDataConfig := gdconf.GetQuestDataById(int32(quest.QuestId))
		if questDataConfig == nil {
			continue
		}
		// TODO 实在不知道客户端要在怎样的情况下 才会发长按10006这个技能 这里先临时改表解决了
		// 是走ability体系计算出来的 操了
		if questDataConfig.QuestId == 35303 {
			questDataConfig.FinishCondList[0].Param[0] = 10067
		}
		for _, questCond := range questDataConfig.FinishCondList {
			if questCond.Type != cond {
				continue
			}
			switch cond {
			case constant.QUEST_FINISH_COND_TYPE_FINISH_PLOT:
				ok := matchParamEqual(questCond.Param, param, 1)
				if !ok {
					continue
				}
				dbQuest.ForceFinishQuest(quest.QuestId)
				updateQuestIdList = append(updateQuestIdList, quest.QuestId)
			case constant.QUEST_FINISH_COND_TYPE_TRIGGER_FIRE:
				// 场景触发器跳了 参数1:触发器id
				ok := matchParamEqual(questCond.Param, param, 1)
				if !ok {
					continue
				}
				dbQuest.ForceFinishQuest(quest.QuestId)
				updateQuestIdList = append(updateQuestIdList, quest.QuestId)
			case constant.QUEST_FINISH_COND_TYPE_UNLOCK_TRANS_POINT:
				// 解锁传送锚点 参数1:场景id 参数2:传送锚点id
				ok := matchParamEqual(questCond.Param, param, 2)
				if !ok {
					continue
				}
				dbQuest.ForceFinishQuest(quest.QuestId)
				updateQuestIdList = append(updateQuestIdList, quest.QuestId)
			case constant.QUEST_FINISH_COND_TYPE_COMPLETE_TALK:
				// 与NPC对话 参数1:对话id
				ok := matchParamEqual(questCond.Param, param, 1)
				if !ok {
					continue
				}
				dbQuest.ForceFinishQuest(quest.QuestId)
				updateQuestIdList = append(updateQuestIdList, quest.QuestId)
			case constant.QUEST_FINISH_COND_TYPE_LUA_NOTIFY:
				// LUA侧通知 复杂参数
				if questCond.ComplexParam != complexParam {
					continue
				}
				dbQuest.ForceFinishQuest(quest.QuestId)
				updateQuestIdList = append(updateQuestIdList, quest.QuestId)
			case constant.QUEST_FINISH_COND_TYPE_SKILL:
				// 使用技能 参数1:技能id
				ok := matchParamEqual(questCond.Param, param, 1)
				if !ok {
					continue
				}
				dbQuest.ForceFinishQuest(quest.QuestId)
				updateQuestIdList = append(updateQuestIdList, quest.QuestId)
			}
		}
	}
	if len(updateQuestIdList) > 0 {
		questList := make([]*proto.Quest, 0)
		for _, questId := range updateQuestIdList {
			pbQuest := g.PacketQuest(player, questId)
			if pbQuest == nil {
				continue
			}
			questList = append(questList, pbQuest)
		}
		g.SendMsg(cmd.QuestListUpdateNotify, player.PlayerID, player.ClientSeq, &proto.QuestListUpdateNotify{
			QuestList: questList,
		})
		parentQuestList := make([]*proto.ParentQuest, 0)
		parentQuestMap := make(map[int32]bool)
		for _, questId := range updateQuestIdList {
			questDataConfig := gdconf.GetQuestDataById(int32(questId))
			if questDataConfig == nil {
				continue
			}
			_, exist := parentQuestMap[questDataConfig.ParentQuestId]
			if exist {
				continue
			}
			parentQuestMap[questDataConfig.ParentQuestId] = true
			finishedParentQuest := true
			subQuestDataMap := gdconf.GetQuestDataMapByParentQuestId(questDataConfig.ParentQuestId)
			for _, subQuestData := range subQuestDataMap {
				quest := dbQuest.GetQuestById(uint32(subQuestData.QuestId))
				if quest == nil {
					finishedParentQuest = false
					break
				}
				if quest.State != constant.QUEST_STATE_FINISHED {
					finishedParentQuest = false
					break
				}
			}
			if finishedParentQuest {
				childQuestList := make([]*proto.ChildQuest, 0)
				for _, subQuestData := range subQuestDataMap {
					childQuestList = append(childQuestList, &proto.ChildQuest{
						State:   constant.QUEST_STATE_FINISHED,
						QuestId: uint32(subQuestData.QuestId),
					})
				}
				parentQuestList = append(parentQuestList, &proto.ParentQuest{
					ParentQuestId:    uint32(questDataConfig.ParentQuestId),
					ParentQuestState: 1,
					IsFinished:       true,
					ChildQuestList:   childQuestList,
				})
			}
		}
		if len(parentQuestList) > 0 {
			g.SendMsg(cmd.FinishedParentQuestUpdateNotify, player.PlayerID, player.ClientSeq, &proto.FinishedParentQuestUpdateNotify{
				ParentQuestList: parentQuestList,
			})
		}
		g.AcceptQuest(player, true)
	}
}

// PacketQuest 打包一个任务
func (g *Game) PacketQuest(player *model.Player, questId uint32) *proto.Quest {
	dbQuest := player.GetDbQuest()
	questDataConfig := gdconf.GetQuestDataById(int32(questId))
	if questDataConfig == nil {
		logger.Error("get quest data config is nil, questId: %v", questId)
		return nil
	}
	quest := dbQuest.GetQuestById(questId)
	if quest == nil {
		logger.Error("get quest is nil, questId: %v", quest.QuestId)
		return nil
	}
	pbQuest := &proto.Quest{
		QuestId:            quest.QuestId,
		State:              uint32(quest.State),
		StartTime:          quest.StartTime,
		ParentQuestId:      uint32(questDataConfig.ParentQuestId),
		StartGameTime:      0,
		AcceptTime:         quest.AcceptTime,
		FinishProgressList: quest.FinishProgressList,
	}
	return pbQuest
}

// PacketQuestListNotify 打包任务列表通知
func (g *Game) PacketQuestListNotify(player *model.Player) *proto.QuestListNotify {
	ntf := &proto.QuestListNotify{
		QuestList: make([]*proto.Quest, 0),
	}
	dbQuest := player.GetDbQuest()
	for _, quest := range dbQuest.GetQuestMap() {
		pbQuest := g.PacketQuest(player, quest.QuestId)
		if pbQuest == nil {
			continue
		}
		ntf.QuestList = append(ntf.QuestList, pbQuest)
	}
	return ntf
}

// PacketFinishedParentQuestNotify 打包已完成父任务列表通知
func (g *Game) PacketFinishedParentQuestNotify(player *model.Player) *proto.FinishedParentQuestNotify {
	ntf := &proto.FinishedParentQuestNotify{
		ParentQuestList: make([]*proto.ParentQuest, 0),
	}
	dbQuest := player.GetDbQuest()
	parentQuestMap := make(map[int32]bool)
	for questId := range dbQuest.GetQuestMap() {
		questDataConfig := gdconf.GetQuestDataById(int32(questId))
		if questDataConfig == nil {
			continue
		}
		_, exist := parentQuestMap[questDataConfig.ParentQuestId]
		if exist {
			continue
		}
		parentQuestMap[questDataConfig.ParentQuestId] = true
		finishedParentQuest := true
		subQuestDataMap := gdconf.GetQuestDataMapByParentQuestId(questDataConfig.ParentQuestId)
		for _, subQuestData := range subQuestDataMap {
			quest := dbQuest.GetQuestById(uint32(subQuestData.QuestId))
			if quest == nil {
				finishedParentQuest = false
				break
			}
			if quest.State != constant.QUEST_STATE_FINISHED {
				finishedParentQuest = false
				break
			}
		}
		if finishedParentQuest {
			childQuestList := make([]*proto.ChildQuest, 0)
			for _, subQuestData := range subQuestDataMap {
				childQuestList = append(childQuestList, &proto.ChildQuest{
					State:   constant.QUEST_STATE_FINISHED,
					QuestId: uint32(subQuestData.QuestId),
				})
			}
			ntf.ParentQuestList = append(ntf.ParentQuestList, &proto.ParentQuest{
				ParentQuestId:    uint32(questDataConfig.ParentQuestId),
				ParentQuestState: 1,
				IsFinished:       true,
				ChildQuestList:   childQuestList,
			})
		}
	}
	return ntf
}

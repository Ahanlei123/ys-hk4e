package game

import (
	"strconv"

	"hk4e/common/constant"
	"hk4e/gdconf"
	"hk4e/gs/model"
	"hk4e/pkg/logger"
	"hk4e/pkg/random"
	"hk4e/protocol/cmd"
	"hk4e/protocol/proto"

	pb "google.golang.org/protobuf/proto"
)

// 大地图模块 大世界相关的所有逻辑

func (g *Game) SceneTransToPointReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.SceneTransToPointReq)
	if player.SceneLoadState != model.SceneEnterDone {
		g.SendError(cmd.SceneTransToPointRsp, player, &proto.SceneTransToPointRsp{}, proto.Retcode_RET_IN_TRANSFER)
		return
	}
	world := WORLD_MANAGER.GetWorldByID(player.WorldId)
	if world == nil {
		g.SendError(cmd.SceneTransToPointRsp, player, &proto.SceneTransToPointRsp{})
		return
	}
	owner := world.GetOwner()
	dbWorld := owner.GetDbWorld()
	dbScene := dbWorld.GetSceneById(req.SceneId)
	if dbScene == nil {
		g.SendError(cmd.SceneTransToPointRsp, player, &proto.SceneTransToPointRsp{}, proto.Retcode_RET_POINT_NOT_UNLOCKED)
		return
	}
	unlock := dbScene.CheckPointUnlock(req.PointId)
	if !unlock {
		// g.SendError(cmd.SceneTransToPointRsp, player, &proto.SceneTransToPointRsp{}, proto.Retcode_RET_POINT_NOT_UNLOCKED)
		// return
	}
	pointDataConfig := gdconf.GetScenePointBySceneIdAndPointId(int32(req.SceneId), int32(req.PointId))
	if pointDataConfig == nil {
		g.SendError(cmd.SceneTransToPointRsp, player, &proto.SceneTransToPointRsp{}, proto.Retcode_RET_POINT_NOT_UNLOCKED)
		return
	}

	// 传送玩家
	g.TeleportPlayer(
		player,
		proto.EnterReason_ENTER_REASON_TRANS_POINT,
		req.SceneId,
		&model.Vector{X: pointDataConfig.TranPos.X, Y: pointDataConfig.TranPos.Y, Z: pointDataConfig.TranPos.Z},
		&model.Vector{X: pointDataConfig.TranRot.X, Y: pointDataConfig.TranRot.Y, Z: pointDataConfig.TranRot.Z},
		0,
		0,
	)

	sceneTransToPointRsp := &proto.SceneTransToPointRsp{
		PointId: req.PointId,
		SceneId: req.SceneId,
	}
	g.SendMsg(cmd.SceneTransToPointRsp, player.PlayerID, player.ClientSeq, sceneTransToPointRsp)
}

func (g *Game) UnlockTransPointReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.UnlockTransPointReq)

	world := WORLD_MANAGER.GetWorldByID(player.WorldId)
	if world == nil {
		g.SendError(cmd.UnlockTransPointRsp, player, &proto.UnlockTransPointRsp{})
		return
	}
	owner := world.GetOwner()
	dbWorld := owner.GetDbWorld()
	dbScene := dbWorld.GetSceneById(req.SceneId)
	if dbScene == nil {
		g.SendError(cmd.UnlockTransPointRsp, player, &proto.UnlockTransPointRsp{}, proto.Retcode_RET_POINT_NOT_UNLOCKED)
		return
	}
	unlock := dbScene.CheckPointUnlock(req.PointId)
	if unlock {
		g.SendError(cmd.UnlockTransPointRsp, player, &proto.UnlockTransPointRsp{}, proto.Retcode_RET_POINT_ALREAY_UNLOCKED)
		return
	}
	dbScene.UnlockPoint(req.PointId)

	g.TriggerQuest(player, constant.QUEST_FINISH_COND_TYPE_UNLOCK_TRANS_POINT, "", int32(req.SceneId), int32(req.PointId))

	g.SendMsg(cmd.ScenePointUnlockNotify, player.PlayerID, player.ClientSeq, &proto.ScenePointUnlockNotify{
		SceneId:         req.SceneId,
		PointList:       []uint32{req.PointId},
		UnhidePointList: nil,
	})
	g.SendSucc(cmd.UnlockTransPointRsp, player, &proto.UnlockTransPointRsp{})
}

func (g *Game) GetScenePointReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.GetScenePointReq)

	world := WORLD_MANAGER.GetWorldByID(player.WorldId)
	if world == nil {
		g.SendError(cmd.GetScenePointRsp, player, &proto.GetScenePointRsp{})
		return
	}
	owner := world.GetOwner()
	dbWorld := owner.GetDbWorld()
	dbScene := dbWorld.GetSceneById(req.SceneId)
	if dbScene == nil {
		g.SendError(cmd.GetScenePointRsp, player, &proto.GetScenePointRsp{})
		return
	}
	getScenePointRsp := &proto.GetScenePointRsp{
		SceneId: req.SceneId,
	}
	areaIdMap := make(map[uint32]bool)
	for _, worldAreaData := range gdconf.GetWorldAreaDataMap() {
		if uint32(worldAreaData.SceneId) == req.SceneId {
			areaIdMap[uint32(worldAreaData.AreaId1)] = true
		}
	}
	areaList := make([]uint32, 0)
	for areaId := range areaIdMap {
		areaList = append(areaList, areaId)
	}
	getScenePointRsp.UnlockAreaList = areaList
	for _, pointId := range dbScene.GetUnlockPointList() {
		pointData := gdconf.GetScenePointBySceneIdAndPointId(int32(req.SceneId), int32(pointId))
		if pointData.IsModelHidden {
			getScenePointRsp.HidePointList = append(getScenePointRsp.HidePointList, pointId)
		}
		getScenePointRsp.UnlockedPointList = append(getScenePointRsp.UnlockedPointList, pointId)
	}
	g.SendMsg(cmd.GetScenePointRsp, player.PlayerID, player.ClientSeq, getScenePointRsp)
}

func (g *Game) MarkMapReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.MarkMapReq)
	if req.Op == proto.MarkMapReq_ADD {
		logger.Debug("user mark type: %v", req.Mark.PointType)
		// 地图标点传送
		if req.Mark.PointType == proto.MapMarkPointType_NPC {
			posYInt, err := strconv.ParseInt(req.Mark.Name, 10, 64)
			if err != nil {
				logger.Error("parse pos y error: %v", err)
				posYInt = 300
			}
			// 传送玩家
			g.TeleportPlayer(
				player,
				proto.EnterReason_ENTER_REASON_GM,
				req.Mark.SceneId,
				&model.Vector{X: float64(req.Mark.Pos.X), Y: float64(posYInt), Z: float64(req.Mark.Pos.Z)},
				new(model.Vector),
				0,
				0,
			)
		}
	}
	g.SendMsg(cmd.MarkMapRsp, player.PlayerID, player.ClientSeq, &proto.MarkMapRsp{})
}

func (g *Game) GetSceneAreaReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.GetSceneAreaReq)

	getSceneAreaRsp := &proto.GetSceneAreaRsp{
		SceneId: req.SceneId,
	}
	areaIdMap := make(map[uint32]bool)
	for _, worldAreaData := range gdconf.GetWorldAreaDataMap() {
		if uint32(worldAreaData.SceneId) == req.SceneId {
			areaIdMap[uint32(worldAreaData.AreaId1)] = true
		}
	}
	areaList := make([]uint32, 0)
	for areaId := range areaIdMap {
		areaList = append(areaList, areaId)
	}
	getSceneAreaRsp.AreaIdList = areaList
	if req.SceneId == 3 {
		getSceneAreaRsp.CityInfoList = []*proto.CityInfo{
			{CityId: 1, Level: 10},
			{CityId: 2, Level: 10},
			{CityId: 3, Level: 10},
			{CityId: 4, Level: 10},
			{CityId: 99, Level: 1},
			{CityId: 100, Level: 1},
			{CityId: 101, Level: 1},
			{CityId: 102, Level: 1},
		}
	}
	g.SendMsg(cmd.GetSceneAreaRsp, player.PlayerID, player.ClientSeq, getSceneAreaRsp)
}

func (g *Game) EnterWorldAreaReq(player *model.Player, payloadMsg pb.Message) {
	logger.Debug("player enter world area, uid: %v", player.PlayerID)
	req := payloadMsg.(*proto.EnterWorldAreaReq)

	logger.Debug("EnterWorldAreaReq: %v", req)

	enterWorldAreaRsp := &proto.EnterWorldAreaRsp{
		AreaType: req.AreaType,
		AreaId:   req.AreaId,
	}
	g.SendMsg(cmd.EnterWorldAreaRsp, player.PlayerID, player.ClientSeq, enterWorldAreaRsp)
}

func (g *Game) ChangeGameTimeReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.ChangeGameTimeReq)
	gameTime := req.GameTime
	world := WORLD_MANAGER.GetWorldByID(player.WorldId)
	if world == nil {
		logger.Error("get world is nil, worldId: %v, uid: %v", player.WorldId, player.PlayerID)
		return
	}
	scene := world.GetSceneById(player.SceneId)
	scene.ChangeGameTime(gameTime)

	for _, scenePlayer := range scene.GetAllPlayer() {
		playerGameTimeNotify := &proto.PlayerGameTimeNotify{
			GameTime: scene.GetGameTime(),
			Uid:      scenePlayer.PlayerID,
		}
		g.SendMsg(cmd.PlayerGameTimeNotify, scenePlayer.PlayerID, scenePlayer.ClientSeq, playerGameTimeNotify)
	}

	changeGameTimeRsp := &proto.ChangeGameTimeRsp{
		CurGameTime: scene.GetGameTime(),
	}
	g.SendMsg(cmd.ChangeGameTimeRsp, player.PlayerID, player.ClientSeq, changeGameTimeRsp)
}

func (g *Game) NpcTalkReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.NpcTalkReq)
	g.TriggerQuest(player, constant.QUEST_FINISH_COND_TYPE_COMPLETE_TALK, "", int32(req.TalkId))
	rsp := &proto.NpcTalkRsp{
		CurTalkId:   req.TalkId,
		NpcEntityId: req.NpcEntityId,
		EntityId:    req.EntityId,
	}
	g.SendMsg(cmd.NpcTalkRsp, player.PlayerID, player.ClientSeq, rsp)
}

func (g *Game) DungeonEntryInfoReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.DungeonEntryInfoReq)
	pointDataConfig := gdconf.GetScenePointBySceneIdAndPointId(int32(req.SceneId), int32(req.PointId))
	if pointDataConfig == nil {
		g.SendError(cmd.DungeonEntryInfoRsp, player, &proto.DungeonEntryInfoRsp{})
		return
	}

	rsp := &proto.DungeonEntryInfoRsp{
		DungeonEntryList: make([]*proto.DungeonEntryInfo, 0),
		PointId:          req.PointId,
	}
	for _, dungeonId := range pointDataConfig.DungeonIds {
		rsp.DungeonEntryList = append(rsp.DungeonEntryList, &proto.DungeonEntryInfo{
			DungeonId: uint32(dungeonId),
		})
	}
	g.SendMsg(cmd.DungeonEntryInfoRsp, player.PlayerID, player.ClientSeq, rsp)
}

func (g *Game) PlayerEnterDungeonReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.PlayerEnterDungeonReq)
	dungeonDataConfig := gdconf.GetDungeonDataById(int32(req.DungeonId))
	if dungeonDataConfig == nil {
		logger.Error("get dungeon data config is nil, dungeonId: %v, uid: %v", req.DungeonId, player.PlayerID)
		return
	}
	sceneLuaConfig := gdconf.GetSceneLuaConfigById(dungeonDataConfig.SceneId)
	if sceneLuaConfig == nil {
		logger.Error("get scene lua config is nil, sceneId: %v, uid: %v", dungeonDataConfig.SceneId, player.PlayerID)
		return
	}
	sceneConfig := sceneLuaConfig.SceneConfig
	g.TeleportPlayer(
		player,
		proto.EnterReason_ENTER_REASON_DUNGEON_ENTER,
		uint32(dungeonDataConfig.SceneId),
		&model.Vector{X: float64(sceneConfig.BornPos.X), Y: float64(sceneConfig.BornPos.Y), Z: float64(sceneConfig.BornPos.Z)},
		&model.Vector{X: float64(sceneConfig.BornRot.X), Y: float64(sceneConfig.BornRot.Y), Z: float64(sceneConfig.BornRot.Z)},
		req.DungeonId,
		req.PointId,
	)

	rsp := &proto.PlayerEnterDungeonRsp{
		DungeonId: req.DungeonId,
		PointId:   req.PointId,
	}
	g.SendMsg(cmd.PlayerEnterDungeonRsp, player.PlayerID, player.ClientSeq, rsp)
}

func (g *Game) PlayerQuitDungeonReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.PlayerQuitDungeonReq)
	world := WORLD_MANAGER.GetWorldByID(player.WorldId)
	if world == nil {
		logger.Error("get world is nil, worldId: %v, uid: %v", player.WorldId, player.PlayerID)
		return
	}
	ctx := world.GetLastEnterSceneContextByUid(player.PlayerID)
	if ctx == nil {
		return
	}
	pointDataConfig := gdconf.GetScenePointBySceneIdAndPointId(int32(ctx.OldSceneId), int32(ctx.OldDungeonPointId))
	if pointDataConfig == nil {
		return
	}
	g.TeleportPlayer(
		player,
		proto.EnterReason_ENTER_REASON_DUNGEON_QUIT,
		ctx.OldSceneId,
		&model.Vector{X: pointDataConfig.TranPos.X, Y: pointDataConfig.TranPos.Y, Z: pointDataConfig.TranPos.Z},
		&model.Vector{X: pointDataConfig.TranRot.X, Y: pointDataConfig.TranRot.Y, Z: pointDataConfig.TranRot.Z},
		0,
		0,
	)

	rsp := &proto.PlayerQuitDungeonRsp{
		PointId: req.PointId,
	}
	g.SendMsg(cmd.PlayerQuitDungeonRsp, player.PlayerID, player.ClientSeq, rsp)
}

func (g *Game) GadgetInteractReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.GadgetInteractReq)
	world := WORLD_MANAGER.GetWorldByID(player.WorldId)
	if world == nil {
		logger.Error("get world is nil, worldId: %v, uid: %v", player.WorldId, player.PlayerID)
		return
	}
	scene := world.GetSceneById(player.SceneId)
	entity := scene.GetEntity(req.GadgetEntityId)
	if entity == nil {
		logger.Error("get entity is nil, entityId: %v, uid: %v", req.GadgetEntityId, player.PlayerID)
		return
	}
	if entity.GetEntityType() != constant.ENTITY_TYPE_GADGET {
		logger.Error("entity type is not gadget, entityType: %v, uid: %v", entity.GetEntityType(), player.PlayerID)
		return
	}
	gadgetEntity := entity.GetGadgetEntity()
	gadgetDataConfig := gdconf.GetGadgetDataById(int32(gadgetEntity.GetGadgetId()))
	if gadgetDataConfig == nil {
		logger.Error("get gadget data config is nil, gadgetId: %v, uid: %v", gadgetEntity.GetGadgetId(), player.PlayerID)
		return
	}
	logger.Debug("[GadgetInteractReq] GadgetData: %+v, EntityId: %v, uid: %v", gadgetDataConfig, entity.GetId(), player.PlayerID)
	interactType := proto.InteractType_INTERACT_NONE
	switch gadgetDataConfig.Type {
	case constant.GADGET_TYPE_GADGET:
		// 掉落物捡起
		interactType = proto.InteractType_INTERACT_PICK_ITEM
		gadgetNormalEntity := gadgetEntity.GetGadgetNormalEntity()
		g.AddUserItem(player.PlayerID, []*ChangeItem{{
			ItemId:      gadgetNormalEntity.GetItemId(),
			ChangeCount: 1,
		}}, true, 0)
		g.KillEntity(player, scene, entity.GetId(), proto.PlayerDieType_PLAYER_DIE_NONE)
	case constant.GADGET_TYPE_ENERGY_BALL:
		// TODO 元素能量球吸收
		interactType = proto.InteractType_INTERACT_PICK_ITEM
	case constant.GADGET_TYPE_GATHER_OBJECT:
		// 采集物摘取
		interactType = proto.InteractType_INTERACT_GATHER
		gadgetNormalEntity := gadgetEntity.GetGadgetNormalEntity()
		g.AddUserItem(player.PlayerID, []*ChangeItem{{
			ItemId:      gadgetNormalEntity.GetItemId(),
			ChangeCount: 1,
		}}, true, 0)
		g.KillEntity(player, scene, entity.GetId(), proto.PlayerDieType_PLAYER_DIE_NONE)
	case constant.GADGET_TYPE_CHEST:
		// 宝箱开启
		interactType = proto.InteractType_INTERACT_OPEN_CHEST
		// 宝箱交互结束 开启宝箱
		if req.OpType == proto.InterOpType_INTER_OP_FINISH {
			// 随机掉落
			g.chestDrop(player, entity)
			// 更新宝箱状态
			g.SendMsg(cmd.WorldChestOpenNotify, player.PlayerID, player.ClientSeq, &proto.WorldChestOpenNotify{
				GroupId:  entity.GetGroupId(),
				SceneId:  scene.GetId(),
				ConfigId: entity.GetConfigId(),
			})
			g.ChangeGadgetState(player, entity.GetId(), constant.GADGET_STATE_CHEST_OPENED)
			g.KillEntity(player, scene, entity.GetId(), proto.PlayerDieType_PLAYER_DIE_NONE)
		}
	}

	rsp := &proto.GadgetInteractRsp{
		GadgetEntityId: req.GadgetEntityId,
		GadgetId:       req.GadgetId,
		OpType:         req.OpType,
		InteractType:   interactType,
	}
	g.SendMsg(cmd.GadgetInteractRsp, player.PlayerID, player.ClientSeq, rsp)
}

func (g *Game) monsterDrop(player *model.Player, entity *Entity) {
	sceneGroupConfig := gdconf.GetSceneGroup(int32(entity.GetGroupId()))
	if sceneGroupConfig == nil {
		logger.Error("get scene group config is nil, groupId: %v, uid: %v", entity.GetGroupId(), player.PlayerID)
		return
	}
	monsterConfig := sceneGroupConfig.MonsterMap[int32(entity.GetConfigId())]
	monsterDropDataConfig := gdconf.GetMonsterDropDataByDropTagAndLevel(monsterConfig.DropTag, monsterConfig.Level)
	if monsterDropDataConfig == nil {
		logger.Error("get monster drop data config is nil, monsterConfig: %v, uid: %v", monsterConfig, player.PlayerID)
		return
	}
	dropDataConfig := gdconf.GetDropDataById(monsterDropDataConfig.DropId)
	if dropDataConfig == nil {
		logger.Error("get drop data config is nil, dropId: %v, uid: %v", monsterDropDataConfig.DropId, player.PlayerID)
		return
	}
	totalItemMap := g.doRandDropFullTimes(dropDataConfig, int(monsterDropDataConfig.DropCount))
	for itemId, count := range totalItemMap {
		itemDataConfig := gdconf.GetItemDataById(int32(itemId))
		if itemDataConfig == nil {
			logger.Error("get item data config is nil, itemId: %v, uid: %v", itemId, player.PlayerID)
			continue
		}
		g.CreateDropGadget(player, entity.pos, uint32(itemDataConfig.GadgetId), itemId, count)
	}
}

func (g *Game) chestDrop(player *model.Player, entity *Entity) {
	sceneGroupConfig := gdconf.GetSceneGroup(int32(entity.GetGroupId()))
	if sceneGroupConfig == nil {
		logger.Error("get scene group config is nil, groupId: %v, uid: %v", entity.GetGroupId(), player.PlayerID)
		return
	}
	gadgetConfig := sceneGroupConfig.GadgetMap[int32(entity.GetConfigId())]
	chestDropDataConfig := gdconf.GetChestDropDataByDropTagAndLevel(gadgetConfig.DropTag, gadgetConfig.Level)
	if chestDropDataConfig == nil {
		logger.Error("get chest drop data config is nil, gadgetConfig: %v, uid: %v", gadgetConfig, player.PlayerID)
		return
	}
	dropDataConfig := gdconf.GetDropDataById(chestDropDataConfig.DropId)
	if dropDataConfig == nil {
		logger.Error("get drop data config is nil, dropId: %v, uid: %v", chestDropDataConfig.DropId, player.PlayerID)
		return
	}
	totalItemMap := g.doRandDropFullTimes(dropDataConfig, int(chestDropDataConfig.DropCount))
	for itemId, count := range totalItemMap {
		itemDataConfig := gdconf.GetItemDataById(int32(itemId))
		if itemDataConfig == nil {
			logger.Error("get item data config is nil, itemId: %v, uid: %v", itemId, player.PlayerID)
			continue
		}
		g.CreateDropGadget(player, entity.pos, uint32(itemDataConfig.GadgetId), itemId, count)
	}
}

func (g *Game) doRandDropFullTimes(dropDataConfig *gdconf.DropData, times int) map[uint32]uint32 {
	totalItemMap := make(map[uint32]uint32)
	for i := 0; i < times; i++ {
		itemMap := g.doRandDropFull(dropDataConfig)
		if itemMap == nil {
			continue
		}
		for itemId, count := range itemMap {
			totalItemMap[itemId] += count
		}
	}
	return totalItemMap
}

func (g *Game) doRandDropFull(dropDataConfig *gdconf.DropData) map[uint32]uint32 {
	itemMap := make(map[uint32]uint32)
	dropList := make([]*gdconf.DropData, 0)
	dropList = append(dropList, dropDataConfig)
	for i := 0; i < 1000; i++ {
		if len(dropList) == 0 {
			// 掉落结束
			return itemMap
		}
		dropMap := g.doRandDropOnce(dropList[0])
		dropList = dropList[1:]
		for dropId, count := range dropMap {
			// 掉落id优先在掉落表里找 找不到就去道具表里找
			subDropDataConfig := gdconf.GetDropDataById(dropId)
			if subDropDataConfig != nil {
				// 添加子掉落
				dropList = append(dropList, subDropDataConfig)
			} else {
				// 添加道具
				itemMap[uint32(dropId)] += uint32(count)
			}
		}
	}
	logger.Error("drop overtimes, drop config: %v", dropDataConfig)
	return nil
}

func (g *Game) doRandDropOnce(dropDataConfig *gdconf.DropData) map[int32]int32 {
	dropMap := make(map[int32]int32)
	switch dropDataConfig.RandomType {
	case gdconf.RandomTypeChoose:
		// RWS随机
		randNum := random.GetRandomInt32(0, dropDataConfig.SubDropTotalWeight-1)
		sumWeight := int32(0)
		for _, subDrop := range dropDataConfig.SubDropList {
			sumWeight += subDrop.Weight
			if sumWeight > randNum {
				dropMap[subDrop.Id] = random.GetRandomInt32(subDrop.CountRange[0], subDrop.CountRange[1])
				break
			}
		}
	case gdconf.RandomTypeIndep:
		// 独立随机
		randNum := random.GetRandomInt32(0, gdconf.RandomTypeIndepWeight-1)
		for _, subDrop := range dropDataConfig.SubDropList {
			if subDrop.Weight > randNum {
				dropMap[subDrop.Id] += random.GetRandomInt32(subDrop.CountRange[0], subDrop.CountRange[1])
			}
		}
	}
	return dropMap
}

// TeleportPlayer 传送玩家通用接口
func (g *Game) TeleportPlayer(
	player *model.Player, enterReason proto.EnterReason,
	sceneId uint32, pos, rot *model.Vector,
	dungeonId, dungeonPointId uint32,
) {
	world := WORLD_MANAGER.GetWorldByID(player.WorldId)
	if world == nil {
		logger.Error("get world is nil, worldId: %v, uid: %v", player.WorldId, player.PlayerID)
		return
	}
	if WORLD_MANAGER.IsBigWorld(world) && sceneId != 3 {
		logger.Error("big world scene not support now, sceneId: %v, uid: %v", sceneId, player.PlayerID)
		return
	}
	newSceneId := sceneId
	oldSceneId := player.SceneId
	oldPos := &model.Vector{X: player.Pos.X, Y: player.Pos.Y, Z: player.Pos.Z}
	jumpScene := false
	if newSceneId != oldSceneId {
		jumpScene = true
	}
	player.SceneJump = jumpScene
	oldScene := world.GetSceneById(oldSceneId)
	activeAvatarId := world.GetPlayerActiveAvatarId(player)
	g.RemoveSceneEntityNotifyBroadcast(oldScene, proto.VisionType_VISION_REMOVE, []uint32{world.GetPlayerWorldAvatarEntityId(player, activeAvatarId)}, false, 0)

	if WORLD_MANAGER.IsBigWorld(world) {
		bigWorldAoi := world.GetBigWorldAoi()
		bigWorldAoi.RemoveObjectFromGridByPos(int64(player.PlayerID), float32(player.Pos.X), float32(player.Pos.Y), float32(player.Pos.Z))
	}

	if jumpScene {
		delTeamEntityNotify := g.PacketDelTeamEntityNotify(oldScene, player)
		g.SendMsg(cmd.DelTeamEntityNotify, player.PlayerID, player.ClientSeq, delTeamEntityNotify)

		oldScene.RemovePlayer(player)
		player.SceneId = newSceneId
		newScene := world.GetSceneById(newSceneId)
		newScene.AddPlayer(player)
	}
	player.SceneLoadState = model.SceneNone
	player.Pos.X, player.Pos.Y, player.Pos.Z = pos.X, pos.Y, pos.Z
	player.Rot.X, player.Rot.Y, player.Rot.Z = rot.X, rot.Y, rot.Z

	var enterType proto.EnterType
	if jumpScene {
		logger.Debug("player jump scene, scene: %v, pos: %v", player.SceneId, player.Pos)
		enterType = proto.EnterType_ENTER_JUMP
		if enterReason == proto.EnterReason_ENTER_REASON_DUNGEON_ENTER {
			logger.Debug("player tp to dungeon scene, sceneId: %v, pos: %v", player.SceneId, player.Pos)
			enterType = proto.EnterType_ENTER_DUNGEON
		}
	} else {
		logger.Debug("player goto scene, scene: %v, pos: %v", player.SceneId, player.Pos)
		enterType = proto.EnterType_ENTER_GOTO
	}
	playerEnterSceneNotify := g.PacketPlayerEnterSceneNotifyTp(player, enterType, enterReason, oldSceneId, oldPos, dungeonId, dungeonPointId)
	g.SendMsg(cmd.PlayerEnterSceneNotify, player.PlayerID, player.ClientSeq, playerEnterSceneNotify)
}

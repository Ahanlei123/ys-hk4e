package game

import (
	"hk4e/common/constant"
	"hk4e/gdconf"
	"hk4e/gs/model"
	"hk4e/pkg/endec"
	"hk4e/pkg/logger"
	"hk4e/protocol/cmd"
	"hk4e/protocol/proto"

	pb "google.golang.org/protobuf/proto"
)

func (g *Game) ChangeAvatarReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.ChangeAvatarReq)
	targetAvatarGuid := req.Guid
	world := WORLD_MANAGER.GetWorldByID(player.WorldId)
	if world == nil {
		logger.Error("get world is nil, worldId: %v, uid: %v", player.WorldId, player.PlayerID)
		return
	}
	scene := world.GetSceneById(player.SceneId)
	targetAvatar, ok := player.GameObjectGuidMap[targetAvatarGuid].(*model.Avatar)
	if !ok {
		logger.Error("target avatar error, avatarGuid: %v", targetAvatarGuid)
		return
	}
	targetAvatarId := targetAvatar.AvatarId
	oldAvatarId := world.GetPlayerActiveAvatarId(player)
	if targetAvatarId == oldAvatarId {
		logger.Error("can not change to the same avatar, uid: %v, oldAvatarId: %v, targetAvatarId: %v", player.PlayerID, oldAvatarId, targetAvatarId)
		return
	}
	newAvatarIndex := world.GetPlayerAvatarIndexByAvatarId(player, targetAvatarId)
	if newAvatarIndex == -1 {
		logger.Error("can not find the target avatar in team, uid: %v, targetAvatarId: %v", player.PlayerID, targetAvatarId)
		return
	}
	if !world.GetMultiplayer() {
		dbTeam := player.GetDbTeam()
		dbTeam.CurrAvatarIndex = uint8(newAvatarIndex)
	}
	world.SetPlayerAvatarIndex(player, newAvatarIndex)
	oldAvatarEntityId := world.GetPlayerWorldAvatarEntityId(player, oldAvatarId)
	oldAvatarEntity := scene.GetEntity(oldAvatarEntityId)
	if oldAvatarEntity == nil {
		logger.Error("can not find old avatar entity, entity id: %v", oldAvatarEntityId)
		return
	}
	oldAvatarEntity.SetMoveState(uint16(proto.MotionState_MOTION_STANDBY))

	sceneEntityDisappearNotify := &proto.SceneEntityDisappearNotify{
		DisappearType: proto.VisionType_VISION_REPLACE,
		EntityList:    []uint32{oldAvatarEntity.GetId()},
	}
	g.SendToSceneA(scene, cmd.SceneEntityDisappearNotify, player.ClientSeq, sceneEntityDisappearNotify)

	newAvatarId := world.GetPlayerActiveAvatarId(player)
	newAvatarEntity := g.PacketSceneEntityInfoAvatar(scene, player, newAvatarId)
	sceneEntityAppearNotify := &proto.SceneEntityAppearNotify{
		AppearType: proto.VisionType_VISION_REPLACE,
		Param:      oldAvatarEntity.GetId(),
		EntityList: []*proto.SceneEntityInfo{newAvatarEntity},
	}
	g.SendToSceneA(scene, cmd.SceneEntityAppearNotify, player.ClientSeq, sceneEntityAppearNotify)

	changeAvatarRsp := &proto.ChangeAvatarRsp{
		CurGuid: targetAvatarGuid,
	}
	g.SendMsg(cmd.ChangeAvatarRsp, player.PlayerID, player.ClientSeq, changeAvatarRsp)
}

func (g *Game) SetUpAvatarTeamReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.SetUpAvatarTeamReq)
	world := WORLD_MANAGER.GetWorldByID(player.WorldId)
	if world == nil {
		logger.Error("get world is nil, worldId: %v, uid: %v", player.WorldId, player.PlayerID)
		return
	}
	if world.GetMultiplayer() {
		g.SendError(cmd.SetUpAvatarTeamRsp, player, &proto.SetUpAvatarTeamRsp{})
		return
	}
	teamId := req.TeamId
	if teamId <= 0 || teamId >= 5 {
		g.SendError(cmd.SetUpAvatarTeamRsp, player, &proto.SetUpAvatarTeamRsp{})
		return
	}
	avatarGuidList := req.AvatarTeamGuidList
	dbTeam := player.GetDbTeam()
	selfTeam := teamId == uint32(dbTeam.GetActiveTeamId())
	if (selfTeam && len(avatarGuidList) == 0) || len(avatarGuidList) > 4 {
		g.SendError(cmd.SetUpAvatarTeamRsp, player, &proto.SetUpAvatarTeamRsp{})
		return
	}
	avatarIdList := make([]uint32, 0)
	dbAvatar := player.GetDbAvatar()
	for _, avatarGuid := range avatarGuidList {
		for avatarId, avatar := range dbAvatar.AvatarMap {
			if avatarGuid == avatar.Guid {
				avatarIdList = append(avatarIdList, avatarId)
			}
		}
	}
	dbTeam.GetTeamByIndex(uint8(teamId - 1)).SetAvatarIdList(avatarIdList)

	avatarTeamUpdateNotify := &proto.AvatarTeamUpdateNotify{
		AvatarTeamMap: make(map[uint32]*proto.AvatarTeam),
	}
	for teamIndex, team := range dbTeam.TeamList {
		avatarTeam := &proto.AvatarTeam{
			TeamName:       team.Name,
			AvatarGuidList: make([]uint64, 0),
		}
		for _, avatarId := range team.GetAvatarIdList() {
			avatarTeam.AvatarGuidList = append(avatarTeam.AvatarGuidList, dbAvatar.AvatarMap[avatarId].Guid)
		}
		avatarTeamUpdateNotify.AvatarTeamMap[uint32(teamIndex)+1] = avatarTeam
	}
	g.SendMsg(cmd.AvatarTeamUpdateNotify, player.PlayerID, player.ClientSeq, avatarTeamUpdateNotify)

	if selfTeam {
		// player.TeamConfig.UpdateTeam()
		world.SetPlayerLocalTeam(player, avatarIdList)
		world.UpdateMultiplayerTeam()
		world.InitPlayerWorldAvatar(player)

		currAvatarGuid := req.CurAvatarGuid
		currAvatar, ok := player.GameObjectGuidMap[currAvatarGuid].(*model.Avatar)
		if !ok {
			logger.Error("avatar error, avatarGuid: %v", currAvatarGuid)
			return
		}
		currAvatarId := currAvatar.AvatarId
		currAvatarIndex := world.GetPlayerAvatarIndexByAvatarId(player, currAvatarId)
		dbTeam.CurrAvatarIndex = uint8(currAvatarIndex)
		world.SetPlayerAvatarIndex(player, currAvatarIndex)

		sceneTeamUpdateNotify := g.PacketSceneTeamUpdateNotify(world, player)
		g.SendMsg(cmd.SceneTeamUpdateNotify, player.PlayerID, player.ClientSeq, sceneTeamUpdateNotify)
	}

	setUpAvatarTeamRsp := &proto.SetUpAvatarTeamRsp{
		TeamId:             req.TeamId,
		CurAvatarGuid:      req.CurAvatarGuid,
		AvatarTeamGuidList: req.AvatarTeamGuidList,
	}
	g.SendMsg(cmd.SetUpAvatarTeamRsp, player.PlayerID, player.ClientSeq, setUpAvatarTeamRsp)
}

func (g *Game) ChooseCurAvatarTeamReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.ChooseCurAvatarTeamReq)
	teamId := req.TeamId
	world := WORLD_MANAGER.GetWorldByID(player.WorldId)
	if world == nil {
		logger.Error("get world is nil, worldId: %v, uid: %v", player.WorldId, player.PlayerID)
		return
	}
	if world.GetMultiplayer() {
		g.SendError(cmd.ChooseCurAvatarTeamRsp, player, &proto.ChooseCurAvatarTeamRsp{})
		return
	}
	dbTeam := player.GetDbTeam()
	team := dbTeam.GetTeamByIndex(uint8(teamId) - 1)
	if team == nil || len(team.GetAvatarIdList()) == 0 {
		return
	}
	dbTeam.CurrTeamIndex = uint8(teamId) - 1
	dbTeam.CurrAvatarIndex = 0
	// player.TeamConfig.UpdateTeam()
	world.SetPlayerAvatarIndex(player, 0)
	world.SetPlayerLocalTeam(player, team.GetAvatarIdList())
	world.UpdateMultiplayerTeam()
	world.InitPlayerWorldAvatar(player)

	sceneTeamUpdateNotify := g.PacketSceneTeamUpdateNotify(world, player)
	g.SendMsg(cmd.SceneTeamUpdateNotify, player.PlayerID, player.ClientSeq, sceneTeamUpdateNotify)

	chooseCurAvatarTeamRsp := &proto.ChooseCurAvatarTeamRsp{
		CurTeamId: teamId,
	}
	g.SendMsg(cmd.ChooseCurAvatarTeamRsp, player.PlayerID, player.ClientSeq, chooseCurAvatarTeamRsp)
}

func (g *Game) ChangeMpTeamAvatarReq(player *model.Player, payloadMsg pb.Message) {
	req := payloadMsg.(*proto.ChangeMpTeamAvatarReq)
	avatarGuidList := req.AvatarGuidList
	world := WORLD_MANAGER.GetWorldByID(player.WorldId)
	if world == nil {
		logger.Error("get world is nil, worldId: %v, uid: %v", player.WorldId, player.PlayerID)
		return
	}
	if WORLD_MANAGER.IsBigWorld(world) || !world.GetMultiplayer() || len(avatarGuidList) == 0 || len(avatarGuidList) > 4 {
		g.SendError(cmd.ChangeMpTeamAvatarRsp, player, &proto.ChangeMpTeamAvatarRsp{})
		return
	}
	avatarIdList := make([]uint32, 0)
	for _, avatarGuid := range avatarGuidList {
		avatar, ok := player.GameObjectGuidMap[avatarGuid].(*model.Avatar)
		if !ok {
			logger.Error("avatar error, avatarGuid: %v", avatarGuid)
			return
		}
		avatarId := avatar.AvatarId
		avatarIdList = append(avatarIdList, avatarId)
	}
	world.SetPlayerLocalTeam(player, avatarIdList)
	world.UpdateMultiplayerTeam()
	world.InitPlayerWorldAvatar(player)

	currAvatarGuid := req.CurAvatarGuid
	currAvatar, ok := player.GameObjectGuidMap[currAvatarGuid].(*model.Avatar)
	if !ok {
		logger.Error("avatar error, avatarGuid: %v", currAvatarGuid)
		return
	}
	currAvatarId := currAvatar.AvatarId
	newAvatarIndex := world.GetPlayerAvatarIndexByAvatarId(player, currAvatarId)
	world.SetPlayerAvatarIndex(player, newAvatarIndex)

	sceneTeamUpdateNotify := g.PacketSceneTeamUpdateNotify(world, player)
	g.SendToWorldA(world, cmd.SceneTeamUpdateNotify, player.ClientSeq, sceneTeamUpdateNotify)

	changeMpTeamAvatarRsp := &proto.ChangeMpTeamAvatarRsp{
		CurAvatarGuid:  req.CurAvatarGuid,
		AvatarGuidList: req.AvatarGuidList,
	}
	g.SendMsg(cmd.ChangeMpTeamAvatarRsp, player.PlayerID, player.ClientSeq, changeMpTeamAvatarRsp)
}

func (g *Game) PacketSceneTeamUpdateNotify(world *World, player *model.Player) *proto.SceneTeamUpdateNotify {
	sceneTeamUpdateNotify := &proto.SceneTeamUpdateNotify{
		IsInMp: world.GetMultiplayer(),
	}
	empty := new(proto.AbilitySyncStateInfo)
	for _, worldAvatar := range world.GetWorldAvatarList() {
		if WORLD_MANAGER.IsBigWorld(world) && worldAvatar.uid != player.PlayerID {
			continue
		}
		worldPlayer := USER_MANAGER.GetOnlineUser(worldAvatar.GetUid())
		if worldPlayer == nil {
			logger.Error("player is nil, uid: %v", worldAvatar.GetUid())
			continue
		}
		worldPlayerScene := world.GetSceneById(worldPlayer.SceneId)
		worldPlayerDbAvatar := worldPlayer.GetDbAvatar()
		worldPlayerAvatar := worldPlayerDbAvatar.AvatarMap[worldAvatar.GetAvatarId()]
		equipIdList := make([]uint32, 0)
		weapon := worldPlayerAvatar.EquipWeapon
		equipIdList = append(equipIdList, weapon.ItemId)
		for _, reliquary := range worldPlayerAvatar.EquipReliquaryMap {
			equipIdList = append(equipIdList, reliquary.ItemId)
		}
		sceneTeamAvatar := &proto.SceneTeamAvatar{
			PlayerUid:         worldPlayer.PlayerID,
			AvatarGuid:        worldPlayerAvatar.Guid,
			SceneId:           worldPlayer.SceneId,
			EntityId:          world.GetPlayerWorldAvatarEntityId(worldPlayer, worldAvatar.GetAvatarId()),
			SceneEntityInfo:   g.PacketSceneEntityInfoAvatar(worldPlayerScene, worldPlayer, worldAvatar.GetAvatarId()),
			WeaponGuid:        worldPlayerAvatar.EquipWeapon.Guid,
			WeaponEntityId:    world.GetPlayerWorldAvatarWeaponEntityId(worldPlayer, worldAvatar.GetAvatarId()),
			IsPlayerCurAvatar: world.GetPlayerActiveAvatarId(worldPlayer) == worldAvatar.GetAvatarId(),
			IsOnScene:         world.GetPlayerActiveAvatarId(worldPlayer) == worldAvatar.GetAvatarId(),
			AvatarAbilityInfo: &proto.AbilitySyncStateInfo{
				IsInited:           len(worldAvatar.GetAbilityList()) != 0,
				DynamicValueMap:    nil,
				AppliedAbilities:   worldAvatar.GetAbilityList(),
				AppliedModifiers:   worldAvatar.GetModifierList(),
				MixinRecoverInfos:  nil,
				SgvDynamicValueMap: nil,
			},
			WeaponAbilityInfo:   empty,
			AbilityControlBlock: new(proto.AbilityControlBlock),
		}
		if world.GetMultiplayer() {
			sceneTeamAvatar.AvatarInfo = g.PacketAvatarInfo(worldPlayerAvatar)
			sceneTeamAvatar.SceneAvatarInfo = g.PacketSceneAvatarInfo(worldPlayerScene, worldPlayer, worldAvatar.GetAvatarId())
		}
		// 角色的ability控制块
		acb := sceneTeamAvatar.AbilityControlBlock
		abilityId := 0
		// 默认ability
		for _, abilityHashCode := range constant.DEFAULT_ABILITY_HASH_CODE {
			abilityId++
			ae := &proto.AbilityEmbryo{
				AbilityId:               uint32(abilityId),
				AbilityNameHash:         uint32(abilityHashCode),
				AbilityOverrideNameHash: uint32(endec.Hk4eAbilityHashCode("Default")),
			}
			acb.AbilityEmbryoList = append(acb.AbilityEmbryoList, ae)
		}
		// 角色ability
		avatarDataConfig := gdconf.GetAvatarDataById(int32(worldAvatar.GetAvatarId()))
		if avatarDataConfig != nil {
			for _, abilityHashCode := range avatarDataConfig.AbilityHashCodeList {
				abilityId++
				ae := &proto.AbilityEmbryo{
					AbilityId:               uint32(abilityId),
					AbilityNameHash:         uint32(abilityHashCode),
					AbilityOverrideNameHash: uint32(endec.Hk4eAbilityHashCode("Default")),
				}
				acb.AbilityEmbryoList = append(acb.AbilityEmbryoList, ae)
			}
		}
		// 技能库ability
		skillDepot := gdconf.GetAvatarSkillDepotDataById(int32(worldPlayerAvatar.SkillDepotId))
		if skillDepot != nil && len(skillDepot.AbilityHashCodeList) != 0 {
			for _, abilityHashCode := range skillDepot.AbilityHashCodeList {
				abilityId++
				ae := &proto.AbilityEmbryo{
					AbilityId:               uint32(abilityId),
					AbilityNameHash:         uint32(abilityHashCode),
					AbilityOverrideNameHash: uint32(endec.Hk4eAbilityHashCode("Default")),
				}
				acb.AbilityEmbryoList = append(acb.AbilityEmbryoList, ae)
			}
		}
		// TODO 队伍ability
		// TODO 装备ability
		sceneTeamUpdateNotify.SceneTeamAvatarList = append(sceneTeamUpdateNotify.SceneTeamAvatarList, sceneTeamAvatar)
	}
	return sceneTeamUpdateNotify
}

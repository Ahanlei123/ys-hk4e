package model

import (
	"time"

	"hk4e/common/constant"
	"hk4e/gdconf"
	"hk4e/pkg/logger"
)

type DbAvatar struct {
	AvatarMap        map[uint32]*Avatar // 角色列表
	MainCharAvatarId uint32             // 主角id
}

func (p *Player) GetDbAvatar() *DbAvatar {
	if p.DbAvatar == nil {
		p.DbAvatar = &DbAvatar{
			AvatarMap: make(map[uint32]*Avatar),
		}
	}
	return p.DbAvatar
}

type Avatar struct {
	AvatarId          uint32               // 角色id
	LifeState         uint16               // 存活状态
	Level             uint8                // 等级
	Exp               uint32               // 经验值
	Promote           uint8                // 突破等阶
	Satiation         uint32               // 饱食度
	SatiationPenalty  uint32               // 饱食度溢出
	CurrHP            float64              // 当前生命值
	CurrEnergy        float64              // 当前元素能量值
	FetterList        []uint32             // 资料解锁条目
	SkillLevelMap     map[uint32]uint32    // 技能等级数据
	SkillDepotId      uint32               // 技能库id
	FlyCloak          uint32               // 当前风之翼
	Costume           uint32               // 当前衣装
	BornTime          int64                // 获得时间
	FetterLevel       uint8                // 好感度等级
	FetterExp         uint32               // 好感度经验
	PromoteRewardMap  map[uint32]bool      // 突破奖励 map[突破等级]是否已被领取
	Guid              uint64               `bson:"-" msgpack:"-"`
	EquipGuidMap      map[uint64]uint64    `bson:"-" msgpack:"-"`
	EquipWeapon       *Weapon              `bson:"-" msgpack:"-"`
	EquipReliquaryMap map[uint8]*Reliquary `bson:"-" msgpack:"-"`
	FightPropMap      map[uint32]float32   `bson:"-" msgpack:"-"`
}

func (a *DbAvatar) InitAllAvatar(player *Player) {
	for _, avatar := range a.AvatarMap {
		a.InitAvatar(player, avatar)
	}
}

func (a *DbAvatar) InitAvatar(player *Player, avatar *Avatar) {
	// 角色战斗属性
	a.InitAvatarFightProp(avatar)
	// guid
	avatar.Guid = player.GetNextGameObjectGuid()
	player.GameObjectGuidMap[avatar.Guid] = GameObject(avatar)
	avatar.EquipGuidMap = make(map[uint64]uint64)
	avatar.EquipReliquaryMap = make(map[uint8]*Reliquary)
	a.AvatarMap[avatar.AvatarId] = avatar
	return
}

// InitAvatarFightProp 初始化角色面板
func (a *DbAvatar) InitAvatarFightProp(avatar *Avatar) {
	avatarDataConfig := gdconf.GetAvatarDataById(int32(avatar.AvatarId))
	if avatarDataConfig == nil {
		logger.Error("avatarDataConfig error, avatarId: %v", avatar.AvatarId)
		return
	}
	avatar.FightPropMap = make(map[uint32]float32)
	avatar.FightPropMap[constant.FIGHT_PROP_NONE] = 0.0
	// 白字攻防血
	avatar.FightPropMap[constant.FIGHT_PROP_BASE_ATTACK] = avatarDataConfig.GetBaseAttackByLevel(avatar.Level)
	avatar.FightPropMap[constant.FIGHT_PROP_BASE_DEFENSE] = avatarDataConfig.GetBaseDefenseByLevel(avatar.Level)
	avatar.FightPropMap[constant.FIGHT_PROP_BASE_HP] = avatarDataConfig.GetBaseHpByLevel(avatar.Level)
	// 白字+绿字攻防血
	avatar.FightPropMap[constant.FIGHT_PROP_CUR_ATTACK] = avatarDataConfig.GetBaseAttackByLevel(avatar.Level)
	avatar.FightPropMap[constant.FIGHT_PROP_CUR_DEFENSE] = avatarDataConfig.GetBaseDefenseByLevel(avatar.Level)
	avatar.FightPropMap[constant.FIGHT_PROP_MAX_HP] = avatarDataConfig.GetBaseHpByLevel(avatar.Level)
	// 当前血量
	avatar.FightPropMap[constant.FIGHT_PROP_CUR_HP] = float32(avatar.CurrHP)
	// 双暴
	avatar.FightPropMap[constant.FIGHT_PROP_CRITICAL] = avatarDataConfig.Critical
	avatar.FightPropMap[constant.FIGHT_PROP_CRITICAL_HURT] = avatarDataConfig.CriticalHurt
	// 元素充能
	avatar.FightPropMap[constant.FIGHT_PROP_CHARGE_EFFICIENCY] = 1.0
	a.SetCurrEnergy(avatar, avatar.CurrEnergy, true)
}

func (a *DbAvatar) AddAvatar(player *Player, avatarId uint32) {
	avatarDataConfig := gdconf.GetAvatarDataById(int32(avatarId))
	if avatarDataConfig == nil {
		logger.Error("avatar data config is nil, avatarId: %v", avatarId)
		return
	}
	skillDepotId := int32(0)
	// 主角可以切换属性 技能库要单独设置 这里默认给风元素的技能库
	if avatarId == 10000005 {
		skillDepotId = 504
	} else if avatarId == 10000007 {
		skillDepotId = 704
	} else {
		skillDepotId = avatarDataConfig.SkillDepotId
	}
	avatarSkillDepotDataConfig := gdconf.GetAvatarSkillDepotDataById(skillDepotId)
	if avatarSkillDepotDataConfig == nil {
		logger.Error("avatar skill depot data config is nil, skillDepotId: %v", skillDepotId)
		return
	}
	avatar := &Avatar{
		AvatarId:          avatarId,
		LifeState:         constant.LIFE_STATE_ALIVE,
		Level:             1,
		Exp:               0,
		Promote:           0,
		Satiation:         0,
		SatiationPenalty:  0,
		CurrHP:            0,
		CurrEnergy:        0,
		FetterList:        make([]uint32, 0),
		SkillLevelMap:     make(map[uint32]uint32),
		SkillDepotId:      uint32(skillDepotId),
		FlyCloak:          140001,
		Costume:           0,
		BornTime:          time.Now().Unix(),
		FetterLevel:       1,
		FetterExp:         0,
		Guid:              0,
		EquipGuidMap:      nil,
		EquipWeapon:       nil,
		EquipReliquaryMap: nil,
		FightPropMap:      nil,
		PromoteRewardMap:  make(map[uint32]bool, len(avatarDataConfig.PromoteRewardMap)),
	}

	// 元素爆发1级
	avatar.SkillLevelMap[uint32(avatarSkillDepotDataConfig.EnergySkill)] = 1
	for _, skillId := range avatarSkillDepotDataConfig.Skills {
		// 小技能1级
		avatar.SkillLevelMap[uint32(skillId)] = 1
	}
	avatar.CurrHP = float64(avatarDataConfig.GetBaseHpByLevel(avatar.Level))

	// 角色突破奖励领取状态
	for promoteLevel := range avatarDataConfig.PromoteRewardMap {
		avatar.PromoteRewardMap[promoteLevel] = false
	}

	a.InitAvatar(player, avatar)
	a.AvatarMap[avatarId] = avatar
}

func (a *DbAvatar) SetCurrEnergy(avatar *Avatar, value float64, max bool) {
	var avatarSkillDataConfig *gdconf.AvatarSkillData = nil
	if avatar.AvatarId == 10000005 || avatar.AvatarId == 10000007 {
		avatarSkillDepotDataConfig := gdconf.GetAvatarSkillDepotDataById(int32(avatar.SkillDepotId))
		if avatarSkillDepotDataConfig == nil {
			return
		}
		avatarSkillDataConfig = gdconf.GetAvatarSkillDataById(avatarSkillDepotDataConfig.EnergySkill)
		if avatarSkillDataConfig == nil {
			return
		}
	} else {
		avatarSkillDataConfig = gdconf.GetAvatarEnergySkillConfig(avatar.AvatarId)
	}
	if avatarSkillDataConfig == nil {
		logger.Error("get avatar energy skill is nil, avatarId: %v", avatar.AvatarId)
		return
	}
	fightPropEnergy := constant.ELEMENT_TYPE_FIGHT_PROP_ENERGY_MAP[int(avatarSkillDataConfig.CostElemType)]
	avatar.FightPropMap[uint32(fightPropEnergy.MaxEnergy)] = float32(avatarSkillDataConfig.CostElemVal)
	if max {
		avatar.FightPropMap[uint32(fightPropEnergy.CurEnergy)] = float32(avatarSkillDataConfig.CostElemVal)
	} else {
		avatar.FightPropMap[uint32(fightPropEnergy.CurEnergy)] = float32(value)
	}
}

func (a *DbAvatar) WearReliquary(avatarId uint32, reliquary *Reliquary) {
	avatar := a.AvatarMap[avatarId]
	reliquaryConfig := gdconf.GetItemDataById(int32(reliquary.ItemId))
	if reliquaryConfig == nil {
		logger.Error("reliquary config error, itemId: %v", reliquary.ItemId)
		return
	}
	avatar.EquipReliquaryMap[uint8(reliquaryConfig.ReliquaryType)] = reliquary
	reliquary.AvatarId = avatarId
	avatar.EquipGuidMap[reliquary.Guid] = reliquary.Guid
}

func (a *DbAvatar) TakeOffReliquary(avatarId uint32, reliquary *Reliquary) {
	avatar := a.AvatarMap[avatarId]
	reliquaryConfig := gdconf.GetItemDataById(int32(reliquary.ItemId))
	if reliquaryConfig == nil {
		logger.Error("reliquary config error, itemId: %v", reliquary.ItemId)
		return
	}
	delete(avatar.EquipReliquaryMap, uint8(reliquaryConfig.ReliquaryType))
	reliquary.AvatarId = 0
	delete(avatar.EquipGuidMap, reliquary.Guid)
}

func (a *DbAvatar) WearWeapon(avatarId uint32, weapon *Weapon) {
	avatar := a.AvatarMap[avatarId]
	avatar.EquipWeapon = weapon
	weapon.AvatarId = avatarId
	avatar.EquipGuidMap[weapon.Guid] = weapon.Guid
}

func (a *DbAvatar) TakeOffWeapon(avatarId uint32, weapon *Weapon) {
	avatar := a.AvatarMap[avatarId]
	avatar.EquipWeapon = nil
	weapon.AvatarId = 0
	delete(avatar.EquipGuidMap, weapon.Guid)
}

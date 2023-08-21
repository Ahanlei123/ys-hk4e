-- 基础信息
local base_info = {
	group_id = 133107139
}

--================================================================
-- 
-- 配置
-- 
--================================================================

-- 怪物
monsters = {
	{ config_id = 139001, monster_id = 28050102, pos = { x = -661.686, y = 171.324, z = 885.299 }, rot = { x = 0.000, y = 210.086, z = 0.000 }, level = 32, drop_tag = "魔法生物", area_id = 8 },
	{ config_id = 139002, monster_id = 28050102, pos = { x = -659.335, y = 170.816, z = 878.544 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 32, drop_tag = "魔法生物", area_id = 8 },
	{ config_id = 139003, monster_id = 28050102, pos = { x = -661.734, y = 175.224, z = 864.219 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 32, drop_tag = "魔法生物", area_id = 8 },
	{ config_id = 139004, monster_id = 28050102, pos = { x = -660.575, y = 175.257, z = 890.744 }, rot = { x = 0.000, y = 210.086, z = 0.000 }, level = 32, drop_tag = "魔法生物", area_id = 8 },
	{ config_id = 139005, monster_id = 28050102, pos = { x = -666.519, y = 176.367, z = 876.258 }, rot = { x = 0.000, y = 210.086, z = 0.000 }, level = 32, drop_tag = "魔法生物", area_id = 8 },
	{ config_id = 139006, monster_id = 28040103, pos = { x = -755.083, y = 191.400, z = 769.916 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 32, drop_tag = "水族", area_id = 8 },
	{ config_id = 139007, monster_id = 28040102, pos = { x = -760.374, y = 165.500, z = 799.844 }, rot = { x = 0.000, y = 124.045, z = 0.000 }, level = 32, drop_tag = "水族", area_id = 8 },
	{ config_id = 139008, monster_id = 28040102, pos = { x = -763.485, y = 165.500, z = 809.138 }, rot = { x = 0.000, y = 124.045, z = 0.000 }, level = 32, drop_tag = "水族", area_id = 8 },
	{ config_id = 139009, monster_id = 28040102, pos = { x = -760.343, y = 165.500, z = 864.346 }, rot = { x = 0.000, y = 124.045, z = 0.000 }, level = 32, drop_tag = "水族", area_id = 8 },
	{ config_id = 139010, monster_id = 28040102, pos = { x = -742.381, y = 165.500, z = 908.006 }, rot = { x = 0.000, y = 124.045, z = 0.000 }, level = 32, drop_tag = "水族", area_id = 8 },
	{ config_id = 139011, monster_id = 28010301, pos = { x = -735.605, y = 165.599, z = 888.658 }, rot = { x = 0.000, y = 14.300, z = 0.000 }, level = 32, drop_tag = "采集动物", area_id = 8 },
	{ config_id = 139012, monster_id = 28010301, pos = { x = -767.372, y = 165.804, z = 926.330 }, rot = { x = 0.000, y = 54.500, z = 0.000 }, level = 32, drop_tag = "采集动物", area_id = 8 },
	{ config_id = 139013, monster_id = 28010301, pos = { x = -754.384, y = 165.864, z = 907.855 }, rot = { x = 0.000, y = 33.435, z = 0.000 }, level = 32, drop_tag = "采集动物", area_id = 8 },
	{ config_id = 139014, monster_id = 28010301, pos = { x = -749.403, y = 167.261, z = 985.683 }, rot = { x = 0.000, y = 33.435, z = 0.000 }, level = 32, drop_tag = "采集动物", area_id = 8 }
}

-- NPC
npcs = {
}

-- 装置
gadgets = {
}

-- 区域
regions = {
}

-- 触发器
triggers = {
}

-- 变量
variables = {
}

--================================================================
-- 
-- 初始化配置
-- 
--================================================================

-- 初始化时创建
init_config = {
	suite = 1,
	end_suite = 0,
	rand_suite = false
}

--================================================================
-- 
-- 小组配置
-- 
--================================================================

suites = {
	{
		-- suite_id = 1,
		-- description = ,
		monsters = { 139001, 139002, 139003, 139004, 139005, 139006, 139007, 139008, 139009, 139010, 139011, 139012, 139013, 139014 },
		gadgets = { },
		regions = { },
		triggers = { },
		rand_weight = 100
	}
}

--================================================================
-- 
-- 触发器
-- 
--================================================================
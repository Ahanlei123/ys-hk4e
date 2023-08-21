-- 基础信息
local base_info = {
	group_id = 133003393
}

--================================================================
-- 
-- 配置
-- 
--================================================================

-- 怪物
monsters = {
	{ config_id = 1278, monster_id = 28030403, pos = { x = 2289.138, y = 207.669, z = -1101.105 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 1, drop_tag = "鸟类", pose_id = 2, area_id = 1 },
	{ config_id = 1281, monster_id = 28030403, pos = { x = 2286.994, y = 208.042, z = -1098.736 }, rot = { x = 0.000, y = 113.736, z = 0.000 }, level = 1, drop_tag = "鸟类", pose_id = 2, area_id = 1 },
	{ config_id = 1285, monster_id = 28030403, pos = { x = 2296.199, y = 207.235, z = -1101.727 }, rot = { x = 0.000, y = 254.900, z = 0.000 }, level = 1, drop_tag = "鸟类", pose_id = 2, area_id = 1 }
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
	rand_suite = true
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
		monsters = { 1278, 1281, 1285 },
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
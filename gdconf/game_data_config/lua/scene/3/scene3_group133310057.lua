-- 基础信息
local base_info = {
	group_id = 133310057
}

--================================================================
-- 
-- 配置
-- 
--================================================================

-- 怪物
monsters = {
	{ config_id = 57001, monster_id = 28060301, pos = { x = -2151.200, y = 269.975, z = 4411.808 }, rot = { x = 0.000, y = 87.823, z = 0.000 }, level = 30, drop_tag = "鸟类", pose_id = 101, area_id = 26 },
	{ config_id = 57002, monster_id = 28060301, pos = { x = -2151.541, y = 269.848, z = 4413.281 }, rot = { x = 0.000, y = 72.586, z = 0.000 }, level = 30, drop_tag = "鸟类", pose_id = 101, area_id = 26 }
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
		monsters = { 57001, 57002 },
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
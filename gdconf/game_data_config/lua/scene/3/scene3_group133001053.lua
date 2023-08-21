-- 基础信息
local base_info = {
	group_id = 133001053
}

--================================================================
-- 
-- 配置
-- 
--================================================================

-- 怪物
monsters = {
	{ config_id = 53001, monster_id = 20010801, pos = { x = 1760.153, y = 195.435, z = -1521.632 }, rot = { x = 0.000, y = 20.054, z = 0.000 }, level = 1, drop_id = 1000100, affix = { 1030 }, area_id = 2 },
	{ config_id = 53002, monster_id = 20010801, pos = { x = 1765.205, y = 195.285, z = -1522.472 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 1, drop_id = 1000100, affix = { 1030 }, area_id = 2 }
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
		monsters = { },
		gadgets = { },
		regions = { },
		triggers = { },
		rand_weight = 100
	},
	{
		-- suite_id = 2,
		-- description = ,
		monsters = { 53001 },
		gadgets = { },
		regions = { },
		triggers = { },
		rand_weight = 100
	},
	{
		-- suite_id = 3,
		-- description = ,
		monsters = { 53002 },
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
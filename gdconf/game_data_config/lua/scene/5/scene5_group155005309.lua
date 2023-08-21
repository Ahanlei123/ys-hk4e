-- 基础信息
local base_info = {
	group_id = 155005309
}

--================================================================
-- 
-- 配置
-- 
--================================================================

-- 怪物
monsters = {
}

-- NPC
npcs = {
}

-- 装置
gadgets = {
	{ config_id = 309001, gadget_id = 70220101, pos = { x = 876.015, y = 217.115, z = 20.159 }, rot = { x = 1.373, y = 359.943, z = 356.770 }, level = 36, isOneoff = true, persistent = true, area_id = 200 },
	{ config_id = 309002, gadget_id = 70220102, pos = { x = 872.747, y = 217.222, z = 23.371 }, rot = { x = 2.866, y = 102.217, z = 2.026 }, level = 36, isOneoff = true, persistent = true, area_id = 200 }
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
		gadgets = { 309001, 309002 },
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
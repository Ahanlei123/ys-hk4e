-- 基础信息
local base_info = {
	group_id = 166001421
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
	{ config_id = 421001, gadget_id = 70500000, pos = { x = 869.885, y = 804.163, z = 512.527 }, rot = { x = 0.000, y = 236.011, z = 0.000 }, level = 36, point_type = 2044, area_id = 300 },
	{ config_id = 421002, gadget_id = 70500000, pos = { x = 874.300, y = 808.413, z = 537.632 }, rot = { x = 358.420, y = 0.998, z = 21.302 }, level = 36, point_type = 2044, area_id = 300 },
	{ config_id = 421003, gadget_id = 70500000, pos = { x = 859.231, y = 808.155, z = 559.577 }, rot = { x = 32.687, y = 38.729, z = 344.003 }, level = 36, point_type = 2044, area_id = 300 },
	{ config_id = 421004, gadget_id = 70500000, pos = { x = 850.924, y = 803.959, z = 579.747 }, rot = { x = 15.920, y = 356.741, z = 0.075 }, level = 36, point_type = 2044, area_id = 300 },
	{ config_id = 421005, gadget_id = 70500000, pos = { x = 847.224, y = 811.246, z = 542.319 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 36, point_type = 1002, area_id = 300 },
	{ config_id = 421006, gadget_id = 70500000, pos = { x = 836.569, y = 815.247, z = 549.505 }, rot = { x = 21.093, y = 40.600, z = 346.023 }, level = 36, point_type = 1002, area_id = 300 },
	{ config_id = 421007, gadget_id = 70500000, pos = { x = 838.864, y = 815.248, z = 551.165 }, rot = { x = 0.000, y = 235.667, z = 0.000 }, level = 36, point_type = 1002, area_id = 300 },
	{ config_id = 421008, gadget_id = 70500000, pos = { x = 845.855, y = 811.310, z = 546.885 }, rot = { x = 0.000, y = 92.491, z = 0.000 }, level = 36, point_type = 1002, area_id = 300 }
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
		gadgets = { 421001, 421002, 421003, 421004, 421005, 421006, 421007, 421008 },
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
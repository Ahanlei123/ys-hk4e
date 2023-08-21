-- 基础信息
local base_info = {
	group_id = 133302678
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
	{ config_id = 678001, gadget_id = 70540041, pos = { x = -1023.221, y = 227.909, z = 2563.414 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 27, area_id = 24 },
	{ config_id = 678002, gadget_id = 70500000, pos = { x = -1022.953, y = 228.891, z = 2563.643 }, rot = { x = 346.837, y = 4.866, z = 14.014 }, level = 27, point_type = 2051, owner = 678001, area_id = 24 },
	{ config_id = 678003, gadget_id = 70500000, pos = { x = -1023.140, y = 228.780, z = 2563.434 }, rot = { x = 344.037, y = 7.209, z = 349.121 }, level = 27, point_type = 2051, owner = 678001, area_id = 24 },
	{ config_id = 678004, gadget_id = 70500000, pos = { x = -1022.922, y = 228.532, z = 2563.264 }, rot = { x = 10.767, y = 8.154, z = 14.727 }, level = 27, point_type = 2051, owner = 678001, area_id = 24 },
	{ config_id = 678005, gadget_id = 70500000, pos = { x = -1023.502, y = 228.435, z = 2563.349 }, rot = { x = 359.538, y = 0.302, z = 342.000 }, level = 27, point_type = 2051, owner = 678001, area_id = 24 },
	{ config_id = 678006, gadget_id = 70290002, pos = { x = -1025.660, y = 227.302, z = 2569.950 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 27, area_id = 24 },
	{ config_id = 678007, gadget_id = 70500000, pos = { x = -1026.731, y = 228.920, z = 2570.664 }, rot = { x = 0.000, y = 265.000, z = 0.000 }, level = 27, point_type = 3001, owner = 678006, area_id = 24 },
	{ config_id = 678008, gadget_id = 70500000, pos = { x = -1024.790, y = 229.930, z = 2568.939 }, rot = { x = 0.000, y = 314.000, z = 0.000 }, level = 27, point_type = 3001, owner = 678006, area_id = 24 },
	{ config_id = 678009, gadget_id = 70500000, pos = { x = -1025.617, y = 230.210, z = 2571.331 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 27, point_type = 3001, owner = 678006, area_id = 24 },
	{ config_id = 678010, gadget_id = 70290507, pos = { x = -1008.163, y = 247.606, z = 2464.888 }, rot = { x = 14.238, y = 338.850, z = 313.044 }, level = 27, area_id = 24 },
	{ config_id = 678011, gadget_id = 70500000, pos = { x = -1008.197, y = 227.376, z = 2526.052 }, rot = { x = 3.713, y = 22.608, z = 355.746 }, level = 27, point_type = 2045, area_id = 24 },
	{ config_id = 678012, gadget_id = 70500000, pos = { x = -926.387, y = 207.522, z = 2536.805 }, rot = { x = 11.321, y = 359.031, z = 350.247 }, level = 27, point_type = 2045, area_id = 24 }
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
		gadgets = { 678001, 678002, 678003, 678004, 678005, 678006, 678007, 678008, 678009, 678010, 678011, 678012 },
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
-- 基础信息
local base_info = {
	group_id = 199002107
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
	{ config_id = 107001, npc_id = 30222, pos = { x = 139.241, y = 120.899, z = -183.488 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, area_id = 401 },
	{ config_id = 107002, npc_id = 30223, pos = { x = 118.915, y = 121.552, z = -179.892 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, area_id = 401 },
	{ config_id = 107003, npc_id = 30224, pos = { x = 120.831, y = 121.018, z = -160.258 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, area_id = 401 }
}

-- 装置
gadgets = {
	{ config_id = 107004, gadget_id = 70710797, pos = { x = 139.238, y = 120.956, z = -183.487 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 1, area_id = 401 },
	{ config_id = 107005, gadget_id = 70710797, pos = { x = 118.915, y = 121.608, z = -179.892 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 1, area_id = 401 },
	{ config_id = 107006, gadget_id = 70710797, pos = { x = 120.792, y = 121.071, z = -160.264 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 1, area_id = 401 },
	{ config_id = 107007, gadget_id = 70710794, pos = { x = 139.238, y = 120.906, z = -183.487 }, rot = { x = 274.254, y = 339.134, z = 42.993 }, level = 1, area_id = 401 },
	{ config_id = 107008, gadget_id = 70710794, pos = { x = 118.915, y = 121.558, z = -179.892 }, rot = { x = 270.764, y = 89.273, z = 271.447 }, level = 1, area_id = 401 },
	{ config_id = 107009, gadget_id = 70710794, pos = { x = 120.792, y = 121.021, z = -160.264 }, rot = { x = 272.958, y = 173.480, z = 116.502 }, level = 1, area_id = 401 }
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
		gadgets = { 107004, 107005, 107006, 107007, 107008, 107009 },
		regions = { },
		triggers = { },
		npcs = { 107001, 107002, 107003 },
		rand_weight = 100
	}
}

--================================================================
-- 
-- 触发器
-- 
--================================================================
-- 基础信息
local base_info = {
	group_id = 133309618
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
	{ config_id = 618001, gadget_id = 70500000, pos = { x = -2438.195, y = -42.644, z = 5743.145 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 32, point_type = 2045, area_id = 27 },
	{ config_id = 618002, gadget_id = 70500000, pos = { x = -2431.713, y = -45.180, z = 5736.656 }, rot = { x = 345.293, y = 1.149, z = 351.118 }, level = 32, point_type = 2045, area_id = 27 }
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
		gadgets = { 618001, 618002 },
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
-- 基础信息
local base_info = {
	group_id = 144001020
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
	{ config_id = 20001, gadget_id = 70500000, pos = { x = 752.552, y = 193.171, z = 221.624 }, rot = { x = 0.000, y = 9.845, z = 0.000 }, level = 1, point_type = 2004, area_id = 102 },
	{ config_id = 20002, gadget_id = 70500000, pos = { x = 716.959, y = 120.345, z = 221.388 }, rot = { x = 0.000, y = 287.900, z = 0.000 }, level = 1, point_type = 2001, area_id = 102 },
	{ config_id = 20003, gadget_id = 70500000, pos = { x = 700.856, y = 120.604, z = 254.776 }, rot = { x = 0.000, y = 336.951, z = 0.000 }, level = 1, point_type = 2001, area_id = 102 },
	{ config_id = 20004, gadget_id = 70500000, pos = { x = 679.955, y = 119.406, z = 218.018 }, rot = { x = 0.000, y = 230.769, z = 0.000 }, level = 1, point_type = 2033, area_id = 102 },
	{ config_id = 20005, gadget_id = 70500000, pos = { x = 669.327, y = 118.566, z = 226.437 }, rot = { x = 0.000, y = 308.017, z = 0.000 }, level = 1, point_type = 2033, area_id = 102 },
	{ config_id = 20006, gadget_id = 70500000, pos = { x = 669.176, y = 119.670, z = 242.158 }, rot = { x = 0.000, y = 229.747, z = 0.000 }, level = 1, point_type = 2033, area_id = 102 },
	{ config_id = 20007, gadget_id = 70290008, pos = { x = 699.824, y = 177.218, z = 318.627 }, rot = { x = 0.000, y = 43.892, z = 0.000 }, level = 1, area_id = 102 },
	{ config_id = 20008, gadget_id = 70500000, pos = { x = 699.824, y = 177.218, z = 318.627 }, rot = { x = 0.000, y = 43.892, z = 0.000 }, level = 1, point_type = 3008, owner = 20007, area_id = 102 },
	{ config_id = 20009, gadget_id = 70290008, pos = { x = 700.512, y = 119.771, z = 290.208 }, rot = { x = 0.000, y = 249.251, z = 0.000 }, level = 1, area_id = 102 },
	{ config_id = 20010, gadget_id = 70500000, pos = { x = 700.512, y = 119.771, z = 290.208 }, rot = { x = 0.000, y = 249.251, z = 0.000 }, level = 1, point_type = 3008, owner = 20009, area_id = 102 }
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
		gadgets = { 20001, 20002, 20003, 20004, 20005, 20006, 20007, 20008, 20009, 20010 },
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
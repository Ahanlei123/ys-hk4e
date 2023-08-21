-- 基础信息
local base_info = {
	group_id = 133304403
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
	{ config_id = 403001, gadget_id = 70300118, pos = { x = -1536.678, y = 244.841, z = 2314.384 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 30, area_id = 21 },
	{ config_id = 403002, gadget_id = 70300118, pos = { x = -1553.550, y = 253.510, z = 2362.726 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 30, area_id = 21 },
	{ config_id = 403003, gadget_id = 70300118, pos = { x = -1517.578, y = 239.321, z = 2339.836 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 30, area_id = 21 },
	-- 第一个电池脑袋上
	{ config_id = 403004, gadget_id = 70350085, pos = { x = -1518.874, y = 245.834, z = 2361.903 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 30, area_id = 21 }
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

-- 废弃数据
garbages = {
	gadgets = {
		{ config_id = 403005, gadget_id = 70350085, pos = { x = -1510.174, y = 241.435, z = 2338.448 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 30, area_id = 21 },
		{ config_id = 403006, gadget_id = 70350085, pos = { x = -1531.887, y = 247.460, z = 2303.001 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 30, area_id = 21 },
		{ config_id = 403007, gadget_id = 70350085, pos = { x = -1555.655, y = 259.810, z = 2373.534 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 30, area_id = 21 }
	}
}

--================================================================
-- 
-- 初始化配置
-- 
--================================================================

-- 初始化时创建
init_config = {
	suite = 2,
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
		gadgets = { 403001, 403002, 403003, 403004 },
		regions = { },
		triggers = { },
		rand_weight = 100
	},
	{
		-- suite_id = 2,
		-- description = ,
		monsters = { },
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
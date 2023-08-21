-- 基础信息
local base_info = {
	group_id = 133004039
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
	{ config_id = 39001, gadget_id = 70500000, pos = { x = 2895.004, y = 233.224, z = -608.040 }, rot = { x = 0.000, y = 231.889, z = 0.000 }, level = 20, point_type = 2002, area_id = 4 },
	{ config_id = 39002, gadget_id = 70500000, pos = { x = 2978.466, y = 264.433, z = -570.691 }, rot = { x = 0.000, y = 50.936, z = 0.000 }, level = 20, point_type = 2004, area_id = 4 },
	{ config_id = 39003, gadget_id = 70500000, pos = { x = 2983.082, y = 245.958, z = -706.682 }, rot = { x = 0.000, y = 93.218, z = 0.000 }, level = 1, point_type = 2001, area_id = 4 },
	{ config_id = 39004, gadget_id = 70500000, pos = { x = 2830.507, y = 254.274, z = -672.373 }, rot = { x = 0.000, y = 95.255, z = 0.000 }, level = 20, point_type = 2001, area_id = 4 },
	{ config_id = 39005, gadget_id = 70500000, pos = { x = 3066.336, y = 258.816, z = -618.644 }, rot = { x = 0.000, y = 248.965, z = 0.000 }, level = 1, point_type = 2004, area_id = 4 },
	{ config_id = 39006, gadget_id = 70500000, pos = { x = 2866.918, y = 239.627, z = -566.949 }, rot = { x = 0.000, y = 245.513, z = 0.000 }, level = 20, point_type = 2001, area_id = 4 },
	{ config_id = 39007, gadget_id = 70500000, pos = { x = 2926.964, y = 259.795, z = -594.754 }, rot = { x = 0.000, y = 8.898, z = 0.000 }, level = 20, point_type = 2007, area_id = 4 },
	{ config_id = 39008, gadget_id = 70500000, pos = { x = 3029.836, y = 274.905, z = -521.369 }, rot = { x = 0.000, y = 76.408, z = 0.000 }, level = 1, point_type = 2007, area_id = 4 },
	{ config_id = 39009, gadget_id = 70500000, pos = { x = 2924.206, y = 252.310, z = -650.190 }, rot = { x = 0.000, y = 177.080, z = 0.000 }, level = 20, point_type = 2007, area_id = 4 }
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
		gadgets = { 39001, 39002, 39003, 39004, 39005, 39006, 39007, 39008, 39009 },
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
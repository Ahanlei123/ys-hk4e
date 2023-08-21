-- 基础信息
local base_info = {
	group_id = 166001292
}

--================================================================
-- 
-- 配置
-- 
--================================================================

-- 怪物
monsters = {
	{ config_id = 292001, monster_id = 20010201, pos = { x = 1108.006, y = 713.910, z = 394.977 }, rot = { x = 0.000, y = 300.665, z = 0.000 }, level = 36, drop_tag = "大史莱姆", pose_id = 201, area_id = 300 }
}

-- NPC
npcs = {
}

-- 装置
gadgets = {
	{ config_id = 292002, gadget_id = 70290200, pos = { x = 1101.258, y = 713.439, z = 434.145 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 36, area_id = 300 }
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
		monsters = { 292001 },
		gadgets = { 292002 },
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
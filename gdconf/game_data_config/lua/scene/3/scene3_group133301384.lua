-- 基础信息
local base_info = {
	group_id = 133301384
}

--================================================================
-- 
-- 配置
-- 
--================================================================

-- 怪物
monsters = {
	{ config_id = 384004, monster_id = 28050106, pos = { x = -805.578, y = 11.004, z = 3193.038 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 30, drop_tag = "魔法生物", area_id = 23 },
	{ config_id = 384005, monster_id = 28050106, pos = { x = -806.929, y = 11.032, z = 3191.489 }, rot = { x = 0.000, y = 326.576, z = 0.000 }, level = 30, drop_tag = "魔法生物", area_id = 23 },
	{ config_id = 384008, monster_id = 28050106, pos = { x = -804.107, y = 11.129, z = 3191.888 }, rot = { x = 0.000, y = 326.576, z = 0.000 }, level = 30, drop_tag = "魔法生物", area_id = 23 }
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

-- 废弃数据
garbages = {
	monsters = {
		{ config_id = 384001, monster_id = 28050106, pos = { x = -798.608, y = 13.593, z = 3181.238 }, rot = { x = 0.000, y = 218.816, z = 0.000 }, level = 30, drop_tag = "魔法生物", area_id = 23 },
		{ config_id = 384002, monster_id = 28050106, pos = { x = -818.724, y = 15.581, z = 3179.558 }, rot = { x = 0.000, y = 293.415, z = 0.000 }, level = 30, drop_tag = "魔法生物", area_id = 23 },
		{ config_id = 384003, monster_id = 28050106, pos = { x = -797.886, y = 14.295, z = 3178.893 }, rot = { x = 0.000, y = 326.576, z = 0.000 }, level = 30, drop_tag = "魔法生物", area_id = 23 },
		{ config_id = 384006, monster_id = 28050106, pos = { x = -818.633, y = 15.807, z = 3179.491 }, rot = { x = 0.000, y = 326.576, z = 0.000 }, level = 30, drop_tag = "魔法生物", area_id = 23 },
		{ config_id = 384007, monster_id = 28050106, pos = { x = -797.891, y = 13.936, z = 3180.060 }, rot = { x = 0.000, y = 326.576, z = 0.000 }, level = 30, drop_tag = "魔法生物", area_id = 23 }
	}
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
		monsters = { 384004, 384005, 384008 },
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
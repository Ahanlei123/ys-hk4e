-- 基础信息
local base_info = {
	group_id = 133302422
}

--================================================================
-- 
-- 配置
-- 
--================================================================

-- 怪物
monsters = {
	{ config_id = 422001, monster_id = 28030501, pos = { x = -701.698, y = 200.000, z = 2826.662 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 27, drop_tag = "鸟类", area_id = 24 },
	{ config_id = 422002, monster_id = 28030501, pos = { x = -704.783, y = 200.000, z = 2824.566 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 27, drop_tag = "鸟类", area_id = 24 },
	{ config_id = 422003, monster_id = 28030501, pos = { x = -698.994, y = 200.000, z = 2821.261 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 27, drop_tag = "鸟类", area_id = 24 },
	{ config_id = 422004, monster_id = 28030501, pos = { x = -810.166, y = 200.000, z = 2783.529 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 27, drop_tag = "鸟类", area_id = 24 },
	{ config_id = 422005, monster_id = 28030501, pos = { x = -526.503, y = 200.000, z = 2887.882 }, rot = { x = 0.000, y = 41.268, z = 0.000 }, level = 27, drop_tag = "鸟类", area_id = 24 },
	{ config_id = 422006, monster_id = 28030501, pos = { x = -508.659, y = 200.000, z = 2886.604 }, rot = { x = 0.000, y = 332.885, z = 0.000 }, level = 27, drop_tag = "鸟类", area_id = 24 },
	{ config_id = 422007, monster_id = 28030501, pos = { x = -512.718, y = 200.000, z = 2901.745 }, rot = { x = 0.000, y = 41.268, z = 0.000 }, level = 27, drop_tag = "鸟类", area_id = 24 }
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
		monsters = { 422001, 422002, 422003, 422004, 422005, 422006, 422007 },
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
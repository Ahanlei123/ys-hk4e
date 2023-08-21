-- 基础信息
local base_info = {
	group_id = 133103310
}

--================================================================
-- 
-- 配置
-- 
--================================================================

-- 怪物
monsters = {
	{ config_id = 310001, monster_id = 22010101, pos = { x = 724.931, y = 246.904, z = 1596.646 }, rot = { x = 0.000, y = 21.074, z = 0.000 }, level = 1, drop_tag = "深渊法师", area_id = 6 },
	{ config_id = 310002, monster_id = 21030401, pos = { x = 727.503, y = 247.137, z = 1598.609 }, rot = { x = 0.000, y = 280.370, z = 0.000 }, level = 1, drop_tag = "丘丘萨满", pose_id = 9012, area_id = 6 },
	{ config_id = 310003, monster_id = 21010201, pos = { x = 731.169, y = 248.242, z = 1598.839 }, rot = { x = 0.000, y = 274.362, z = 0.000 }, level = 1, drop_tag = "丘丘人", area_id = 6 },
	{ config_id = 310004, monster_id = 21010201, pos = { x = 726.724, y = 247.018, z = 1603.765 }, rot = { x = 0.000, y = 95.655, z = 0.000 }, level = 1, drop_tag = "丘丘人", area_id = 6 },
	{ config_id = 310005, monster_id = 21010601, pos = { x = 723.925, y = 246.284, z = 1597.897 }, rot = { x = 0.000, y = 78.063, z = 0.000 }, level = 1, drop_tag = "丘丘人", area_id = 6 },
	{ config_id = 310006, monster_id = 21011001, pos = { x = 721.544, y = 249.965, z = 1599.194 }, rot = { x = 0.000, y = 356.285, z = 0.000 }, level = 1, drop_tag = "远程丘丘人", pose_id = 32, area_id = 6 }
}

-- NPC
npcs = {
}

-- 装置
gadgets = {
	{ config_id = 310007, gadget_id = 70300107, pos = { x = 725.670, y = 246.644, z = 1598.404 }, rot = { x = 0.000, y = 304.162, z = 0.000 }, level = 1, state = GadgetState.GearStart, area_id = 6 },
	{ config_id = 310008, gadget_id = 70220013, pos = { x = 723.072, y = 247.166, z = 1592.076 }, rot = { x = 357.832, y = 358.084, z = 355.802 }, level = 1, area_id = 6 },
	{ config_id = 310009, gadget_id = 70220013, pos = { x = 728.293, y = 247.710, z = 1593.776 }, rot = { x = 12.520, y = 260.785, z = 350.750 }, level = 1, area_id = 6 },
	{ config_id = 310010, gadget_id = 70220014, pos = { x = 715.392, y = 245.644, z = 1598.846 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 1, area_id = 6 },
	{ config_id = 310011, gadget_id = 70300086, pos = { x = 725.650, y = 247.523, z = 1593.067 }, rot = { x = 0.000, y = 265.828, z = 0.000 }, level = 1, area_id = 6 },
	{ config_id = 310012, gadget_id = 70220014, pos = { x = 714.497, y = 245.611, z = 1597.152 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 1, area_id = 6 }
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
		-- description = suite_1,
		monsters = { 310001, 310002, 310003, 310004, 310005, 310006 },
		gadgets = { 310007, 310008, 310009, 310010, 310011, 310012 },
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
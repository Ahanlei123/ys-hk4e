-- 基础信息
local base_info = {
	group_id = 133304361
}

--================================================================
-- 
-- 配置
-- 
--================================================================

-- 怪物
monsters = {
	{ config_id = 361002, monster_id = 26120101, pos = { x = -1015.488, y = 241.610, z = 2914.177 }, rot = { x = 0.000, y = 178.016, z = 0.000 }, level = 27, drop_tag = "大蕈兽", disableWander = true, pose_id = 201, area_id = 23 },
	{ config_id = 361003, monster_id = 26090401, pos = { x = -1013.097, y = 239.053, z = 2901.430 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 30, drop_tag = "蕈兽", disableWander = true, pose_id = 102, area_id = 23 },
	{ config_id = 361004, monster_id = 26090401, pos = { x = -1019.817, y = 240.187, z = 2902.214 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 30, drop_tag = "蕈兽", disableWander = true, pose_id = 102, area_id = 23 }
}

-- NPC
npcs = {
}

-- 装置
gadgets = {
	{ config_id = 361001, gadget_id = 70217020, pos = { x = -1025.868, y = 242.752, z = 2919.353 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 26, drop_tag = "摩拉石箱须弥", isOneoff = true, persistent = true, explore = { name = "chest", exp = 4 }, area_id = 23 }
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
		monsters = { 361002, 361003, 361004 },
		gadgets = { 361001 },
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
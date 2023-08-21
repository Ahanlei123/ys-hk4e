-- 基础信息
local base_info = {
	group_id = 133313135
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
	{ config_id = 135001, gadget_id = 70330409, pos = { x = -892.580, y = -247.427, z = 5975.835 }, rot = { x = 0.000, y = 0.000, z = 0.000 }, level = 32, area_id = 32 },
	{ config_id = 135002, gadget_id = 70330409, pos = { x = -864.643, y = -246.921, z = 5988.139 }, rot = { x = 0.000, y = 0.000, z = 11.361 }, level = 32, area_id = 32 },
	{ config_id = 135003, gadget_id = 70330409, pos = { x = -856.838, y = -249.111, z = 6026.092 }, rot = { x = 22.683, y = 87.513, z = 349.485 }, level = 32, area_id = 32 },
	{ config_id = 135004, gadget_id = 70330409, pos = { x = -816.167, y = -256.840, z = 6038.429 }, rot = { x = 353.386, y = 78.625, z = 359.962 }, level = 32, area_id = 32 }
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
		gadgets = { 135001, 135002, 135003, 135004 },
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
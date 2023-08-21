-- 基础信息
local base_info = {
	group_id = 133217267
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
	{ config_id = 267001, gadget_id = 70210101, pos = { x = -4323.125, y = 202.079, z = -3852.812 }, rot = { x = 2.272, y = 291.409, z = 346.969 }, level = 26, drop_tag = "搜刮点解谜矿石稻妻", area_id = 14 },
	{ config_id = 267002, gadget_id = 70210101, pos = { x = -4323.757, y = 200.659, z = -3852.709 }, rot = { x = 358.745, y = 291.669, z = 359.898 }, level = 26, drop_tag = "搜刮点解谜矿石稻妻", area_id = 14 },
	{ config_id = 267003, gadget_id = 70210101, pos = { x = -4323.958, y = 201.305, z = -3849.429 }, rot = { x = 345.749, y = 290.799, z = 14.485 }, level = 26, drop_tag = "搜刮点解谜果蔬稻妻", area_id = 14 }
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
		gadgets = { 267001, 267002, 267003 },
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
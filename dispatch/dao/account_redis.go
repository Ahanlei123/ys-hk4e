package dao

import (
	"context"
)

const RedisPlayerKeyPrefix = "HK4E"

const (
	AccountIdRedisKey          = "AccountId"
	AccountIdBegin      uint32 = 10000
	YuanShenUidRedisKey        = "YuanShenUid"
	YuanShenUidBegin    uint32 = 100000000
)

func (d *Dao) GetNextAccountId() (uint32, error) {
	return d.redisInc(RedisPlayerKeyPrefix + ":" + AccountIdRedisKey)
}

func (d *Dao) GetNextYuanShenUid() (uint32, error) {
	return d.redisInc(RedisPlayerKeyPrefix + ":" + YuanShenUidRedisKey)
}

func (d *Dao) redisInc(keyName string) (uint32, error) {
	var exist int64 = 0
	var err error = nil
	if d.redisCluster != nil {
		exist, err = d.redisCluster.Exists(context.TODO(), keyName).Result()
	} else {
		exist, err = d.redis.Exists(context.TODO(), keyName).Result()
	}
	if err != nil {
		return 0, err
	}
	if exist == 0 {
		var value uint32 = 0
		if keyName == RedisPlayerKeyPrefix+":"+AccountIdRedisKey {
			value = AccountIdBegin
		} else if keyName == RedisPlayerKeyPrefix+":"+YuanShenUidRedisKey {
			value = YuanShenUidBegin
		}
		var err error = nil
		if d.redisCluster != nil {
			err = d.redisCluster.Set(context.TODO(), keyName, value, 0).Err()
		} else {
			err = d.redis.Set(context.TODO(), keyName, value, 0).Err()
		}
		if err != nil {
			return 0, err
		}
	}
	var id int64 = 0
	if d.redisCluster != nil {
		id, err = d.redisCluster.Incr(context.TODO(), keyName).Result()
	} else {
		id, err = d.redis.Incr(context.TODO(), keyName).Result()
	}
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

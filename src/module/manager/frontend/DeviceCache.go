package frontend

import (
	"context"
	"github.com/go-redis/redis"
)

type Redis struct {
	_ctx 	context.Context
	_redis  *redis.Client
}

func (obj * Redis)initRedis() {
	obj._ctx = context.Background()
	obj._redis	= redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
}


type DeviceCache struct {
	redis *Redis
}

func (obj * DeviceCache) Init() {
	//todo:应该度配置文件得到redis的相关配置参数
	obj.redis = new (Redis)
	obj.redis.initRedis()
	obj.redis._redis.
}




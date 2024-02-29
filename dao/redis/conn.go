package redis

import (
	"filestore-server/config"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
)

type RedisClient struct {
	conn *redis.Client
}

var (
	once sync.Once
	rdb  *RedisClient
)

func GetRDBInstance() *RedisClient {
	once.Do(func() {
		rdb = &RedisClient{}
	})
	return rdb
}

func (r *RedisClient) GetConn() *redis.Client {
	return r.conn
}

func (r *RedisClient) Init(config config.RedisConf) {
	r.conn = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.PassWord,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})
	_, err := r.conn.Ping().Result()
	if err != nil {
		panic(err)
	}
}

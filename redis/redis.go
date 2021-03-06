package redis

import (
	"fmt"
	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	"time"
)

var pool *redis.Pool
var redisLock *redsync.Redsync

func InitRedis(host, port, password string) error {
	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := conn.Do("AUTH", password); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, err
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
		MaxIdle:         20,
		MaxActive:       50,
		IdleTimeout:     240 * time.Second,
		Wait:            false,
		MaxConnLifetime: 0,
	}

	redisLock = redsync.New([]redsync.Pool{pool})
	return nil
}

func GetRedisConn() (redis.Conn, error) {
	conn := pool.Get()
	return conn, conn.Err()
}

func GetRedisLock(key string, expireTime time.Duration) *redsync.Mutex {
	return redisLock.NewMutex(key, redsync.SetExpiry(expireTime))
}

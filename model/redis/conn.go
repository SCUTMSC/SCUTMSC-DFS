package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	password = "password"
	ip       = "127.0.0.1"
	port     = "6379"
)

var pool *redis.Pool

func init() {
	var err error

	path := strings.Join([]string{ip, ":", port}, "")
	pool = &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			cache, err := redis.Dial("tcp", path)
			if err != nil {
				fmt.Println("Fail to dial redis, err: \n", err.Error())
				return nil, err
			}

			_, err = cache.Do("AUTH", password)
			if err != nil {
				fmt.Println("Fail to auth redis, err: \n", err.Error())
				return nil, err
			}

			return cache, nil
		},
		TestOnBorrow: func(cache redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err = cache.Do("PING")
			if err != nil {
				fmt.Println("Fail to ping redis, err: \n", err.Error())
				return err
			}
			return nil
		},
	}
}

func RedisPool() *redis.Pool {
	return pool
}

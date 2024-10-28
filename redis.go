package gu

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type RedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int
	Ctx      context.Context
}

var DefaultRedisConfig = RedisConfig{
	Host:     "127.0.0.1",
	Port:     6379,
	Username: "",
	Password: "",
	DB:       0,
}

func MustInitRedisClient(config ...RedisConfig) *redis.Client {
	var (
		err     error
		dConfig RedisConfig
	)

	if len(config) > 0 {
		dConfig = config[0]
	} else {
		dConfig = DefaultRedisConfig
	}

	if dConfig.Ctx == nil {
		dConfig.Ctx = Timeout()
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", dConfig.Host, dConfig.Port),
		Username: dConfig.Username,
		Password: dConfig.Password,
		DB:       dConfig.DB,
	})

	if err = RedisClient.Ping(dConfig.Ctx).Err(); err != nil {
		logrus.Panicf("redis client init err: %v\n", err)
	}

	logrus.Info("successfully init redis client")

	return RedisClient
}

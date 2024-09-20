package redisutil

import (
	"context"
	"crypto/tls"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"

	"log"
)

var RDB *redis.Client

func Init() {
	log.Printf("[INFO] Redis trying connect to tcp://%s/%d", viper.GetString("REDIS_ADDR"), viper.GetInt("REDIS_DB"))

	option := &redis.Options{
		Addr:            viper.GetString("REDIS_ADDR"),
		Password:        viper.GetString("REDIS_PWD"),
		DB:              viper.GetInt("REDIS_DB"),
		MinIdleConns:    5,
		MaxIdleConns:    50,
		MaxRetries:      3,
		MaxActiveConns:  50,
		ConnMaxLifetime: 5 * time.Minute,
	}
	if viper.GetBool("REDIS_TLS") {
		option.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	rdb := redis.NewClient(option)

	RDB = rdb

	// test redis
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("[ERROR] Redis ping failed: %v", err)
	}

	log.Printf("[INFO] Redis init success, pong: %s ", pong)
}

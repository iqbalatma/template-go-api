package config

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func ConnectRDB() {
	db, err := strconv.Atoi(AppConfig.RedisDB)
	if err != nil {
		panic("Invalid Redis DB")
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     AppConfig.RedisHost + ":" + AppConfig.RedisPort,
		Password: AppConfig.RedisPassword,
		DB:       db,
	})

	if err := RDB.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("Cannot connect to Redis: %s", err.Error()))
	}
}

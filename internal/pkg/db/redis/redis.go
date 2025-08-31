package redis

import (
	"strconv"

	"github.com/team-xquare/deployment-platform/internal/pkg/config"

	"github.com/go-redis/redis/v8"
)

func NewConnection() (*redis.Client, error) {
	db, _ := strconv.Atoi(config.AppConfig.RedisDB)

	client := redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisHost + ":" + config.AppConfig.RedisPort,
		Password: config.AppConfig.RedisPassword,
		DB:       db,
	})

	return client, nil
}

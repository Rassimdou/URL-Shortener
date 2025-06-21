package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v8"
)

var Ctx = context.Background()

func CreatrClient(dbNo int) *redis.Client {
	rdb := redix.NewClient(&redis.Options{
		addr:     os.Getenv("DB_ADDR "),
		Password: os.Getenv("DB_PASS"),
		DB:       dbNo,
	})
	return rdb
}

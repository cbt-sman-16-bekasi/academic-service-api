package redisstore

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func init() {
	_ = godotenv.Load()
	db, _ := strconv.Atoi(os.Getenv("yon.redis.db"))
	err := Init(Config{
		Addr:     os.Getenv("yon.redis.host"),
		Password: os.Getenv("yon.redis.password"),
		DB:       db,
	})
	if err != nil {
		panic(err)
	}
}

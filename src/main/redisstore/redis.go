package redisstore

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

// Config untuk inisialisasi
type Config struct {
	Addr     string
	Password string
	DB       int
}

// Init Redis connection with retry
func Init(cfg Config) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	var err error
	for i := 0; i < 3; i++ {
		_, err = rdb.Ping(ctx).Result()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 2)
	}
	return err
}

// Set key dengan TTL
func Set(key string, value interface{}, ttl time.Duration) error {
	return rdb.Set(ctx, key, value, ttl).Err()
}

// Get key
func Get(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

// Delete key
func Del(key string) error {
	return rdb.Del(ctx, key).Err()
}

// Set object sebagai JSON
func SetJSON(key string, data interface{}, ttl time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return Set(key, string(b), ttl)
}

// Get key dan unmarshal ke struct
func GetJSON(key string, out interface{}) error {
	val, err := Get(key)
	if err != nil {
		return err
	}
	if val == "" {
		return errors.New("no data found")
	}
	return json.Unmarshal([]byte(val), out)
}

func DeleteByPrefix(prefix string) error {
	var cursor uint64
	var err error
	var keys []string

	for {
		var scanned []string
		scanned, cursor, err = rdb.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}

		if len(scanned) > 0 {
			keys = append(keys, scanned...)
		}

		if cursor == 0 {
			break
		}
	}

	if len(keys) == 0 {
		return nil
	}

	return rdb.Del(ctx, keys...).Err()
}

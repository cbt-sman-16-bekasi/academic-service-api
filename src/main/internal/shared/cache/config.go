package cache

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/config"
	"github.com/rs/zerolog/log"
)

func init() {
	cfg := config.GetRedis()
	err := Init(Config{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to connect to Redis, caching disabled")
	}
}

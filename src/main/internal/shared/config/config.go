package config

import (
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	once   sync.Once
	config *Config
)

// Config holds all application configuration
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Security SecurityConfig `mapstructure:"security"`
	Minio    MinioConfig    `mapstructure:"minio"`
}

// AppConfig holds application configuration
type AppConfig struct {
	Env         string `mapstructure:"env"`
	Port        string `mapstructure:"port"`
	Name        string `mapstructure:"name"`
	ContextPath string `mapstructure:"contextPath"`
	Timezone    string `mapstructure:"timezone"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type          string     `mapstructure:"type"`
	Host          string     `mapstructure:"host"`
	Port          string     `mapstructure:"port"`
	Name          string     `mapstructure:"name"`
	User          string     `mapstructure:"user"`
	Password      string     `mapstructure:"password"`
	Pool          PoolConfig `mapstructure:"pool"`
	SlowThreshold int        `mapstructure:"slowThreshold"` // in milliseconds
	LogLevel      string     `mapstructure:"logLevel"`      // silent, error, warn, info
}

// PoolConfig holds database connection pool settings
type PoolConfig struct {
	MaxOpen     int `mapstructure:"maxOpen"`     // Max open connections
	MaxIdle     int `mapstructure:"maxIdle"`     // Max idle connections
	MaxLifetime int `mapstructure:"maxLifetime"` // Max lifetime in seconds
	MaxIdleTime int `mapstructure:"maxIdleTime"` // Max idle time in seconds
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	SecretKey string `mapstructure:"secretKey"`
}

// MinioConfig holds MinIO configuration
type MinioConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
	SSL       bool   `mapstructure:"ssl"`
}

// Load initializes and returns the application configuration
func Load() *Config {
	once.Do(func() {
		v := viper.New()

		// Config file settings
		v.SetConfigName("config")
		v.SetConfigType("yml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/academic-service/")

		// Set defaults
		setDefaults(v)

		// Enable environment variable override
		v.AutomaticEnv()

		// Read config file
		if err := v.ReadInConfig(); err != nil {
			log.Warn().Err(err).Msg("Config file not found, using defaults and environment variables")
		} else {
			log.Info().Str("file", v.ConfigFileUsed()).Msg("Config loaded")
		}

		// Unmarshal config
		config = &Config{}
		if err := v.Unmarshal(config); err != nil {
			log.Fatal().Err(err).Msg("Failed to unmarshal config")
		}
	})

	return config
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.env", "development")
	v.SetDefault("app.port", "8080")
	v.SetDefault("app.name", "Academic Service API")
	v.SetDefault("app.contextPath", "/")
	v.SetDefault("app.timezone", "Asia/Jakarta")

	// Database defaults
	v.SetDefault("database.type", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.name", "academic_db")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "")
	v.SetDefault("database.pool.maxOpen", 25)
	v.SetDefault("database.pool.maxIdle", 10)
	v.SetDefault("database.pool.maxLifetime", 300)
	v.SetDefault("database.pool.maxIdleTime", 60)
	v.SetDefault("database.slowThreshold", 200)
	v.SetDefault("database.logLevel", "warn")

	// Redis defaults
	v.SetDefault("redis.host", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	// Security defaults
	v.SetDefault("security.secretKey", "default-secret-key")

	// Minio defaults
	v.SetDefault("minio.endpoint", "localhost:9000")
	v.SetDefault("minio.bucket", "academic-bucket")
	v.SetDefault("minio.accessKey", "minioadmin")
	v.SetDefault("minio.secretKey", "minioadmin")
	v.SetDefault("minio.ssl", false)
}

// Get returns the loaded configuration (Load must be called first)
func Get() *Config {
	if config == nil {
		return Load()
	}
	return config
}

// GetApp returns app configuration
func GetApp() *AppConfig {
	return &Get().App
}

// GetDatabase returns database configuration
func GetDatabase() *DatabaseConfig {
	return &Get().Database
}

// GetRedis returns Redis configuration
func GetRedis() *RedisConfig {
	return &Get().Redis
}

// GetSecurity returns security configuration
func GetSecurity() *SecurityConfig {
	return &Get().Security
}

// GetMinio returns MinIO configuration
func GetMinio() *MinioConfig {
	return &Get().Minio
}

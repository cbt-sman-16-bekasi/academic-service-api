package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/config"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	tableMigration []interface{}
	migrationMutex sync.Mutex
)

// DatabaseResult is the FX result for database provider
type DatabaseResult struct {
	fx.Out

	DB    *gorm.DB
	SqlDB *sql.DB
}

// DatabaseParams is the FX params for database provider
type DatabaseParams struct {
	fx.In

	Config    *config.Config
	Lifecycle fx.Lifecycle
}

// NewDatabase creates a new database connection with FX lifecycle management
func NewDatabase(p DatabaseParams) (DatabaseResult, error) {
	cfg := p.Config.Database

	log.Info().Msg("🏗️ Initializing database connection...")
	log.Info().
		Str("type", cfg.Type).
		Str("host", cfg.Host).
		Str("port", cfg.Port).
		Str("name", cfg.Name).
		Str("user", cfg.User).
		Int("maxOpen", cfg.Pool.MaxOpen).
		Int("maxIdle", cfg.Pool.MaxIdle).
		Int("maxLifetime", cfg.Pool.MaxLifetime).
		Msg("Database configuration")

	// Create GORM logger
	gormLog := newGormLogger(cfg.LogLevel, cfg.SlowThreshold)

	// Get dialector based on database type
	dialector := getDialector(cfg)

	// Open database connection
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger:                 gormLog,
		SkipDefaultTransaction: true, // Better performance
		PrepareStmt:            true, // Cache prepared statements
	})
	if err != nil {
		log.Error().Err(err).Msg("❌ Failed to connect to database")
		return DatabaseResult{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("❌ Failed to get underlying sql.DB")
		return DatabaseResult{}, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Configure connection pool
	configurePool(sqlDB, cfg.Pool)

	// Register lifecycle hooks
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Verify connection on start
			if err := sqlDB.PingContext(ctx); err != nil {
				log.Error().Err(err).Msg("❌ Database ping failed")
				return err
			}
			log.Info().Msg("✅ Database connection verified")

			// Run migrations
			runMigrations(db)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("🔌 Closing database connection...")
			if err := sqlDB.Close(); err != nil {
				log.Error().Err(err).Msg("❌ Failed to close database connection")
				return err
			}
			log.Info().Msg("✅ Database connection closed")
			return nil
		},
	})

	log.Info().Msg("✅ Database initialized successfully")

	return DatabaseResult{
		DB:    db,
		SqlDB: sqlDB,
	}, nil
}

// getDialector returns the appropriate GORM dialector
func getDialector(cfg config.DatabaseConfig) gorm.Dialector {
	switch cfg.Type {
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
		return postgres.Open(dsn)

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
		return mysql.Open(dsn)

	case "sqlite":
		return sqlite.Open(cfg.Name)

	default:
		log.Fatal().Str("type", cfg.Type).Msg("Unsupported database type")
		return nil
	}
}

// configurePool configures the database connection pool
func configurePool(sqlDB *sql.DB, pool config.PoolConfig) {
	sqlDB.SetMaxOpenConns(pool.MaxOpen)
	sqlDB.SetMaxIdleConns(pool.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(pool.MaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(pool.MaxIdleTime) * time.Second)

	log.Info().
		Int("maxOpen", pool.MaxOpen).
		Int("maxIdle", pool.MaxIdle).
		Dur("maxLifetime", time.Duration(pool.MaxLifetime)*time.Second).
		Dur("maxIdleTime", time.Duration(pool.MaxIdleTime)*time.Second).
		Msg("Connection pool configured")
}

// newGormLogger creates a GORM logger based on config
func newGormLogger(level string, slowThreshold int) gormlogger.Interface {
	var logLevel gormlogger.LogLevel
	switch level {
	case "silent":
		logLevel = gormlogger.Silent
	case "error":
		logLevel = gormlogger.Error
	case "info":
		logLevel = gormlogger.Info
	default:
		logLevel = gormlogger.Warn
	}

	return gormlogger.New(
		&zerologWriter{},
		gormlogger.Config{
			SlowThreshold:             time.Duration(slowThreshold) * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}

// zerologWriter adapts zerolog for GORM logger interface
type zerologWriter struct{}

func (z *zerologWriter) Printf(format string, args ...interface{}) {
	log.Debug().Msgf(format, args...)
}

// runMigrations runs auto migrations for registered tables
func runMigrations(db *gorm.DB) {
	if len(tableMigration) == 0 {
		log.Info().Msg("No migrations registered")
		return
	}

	log.Info().Int("tables", len(tableMigration)).Msg("🚀 Running migrations...")
	if err := db.AutoMigrate(tableMigration...); err != nil {
		log.Error().Err(err).Msg("❌ Migration failed")
		return
	}
	log.Info().Msg("✅ Migrations completed successfully")
}

// MigrationRegister registers entities for auto migration
func MigrationRegister(tables ...interface{}) {
	migrationMutex.Lock()
	defer migrationMutex.Unlock()
	tableMigration = append(tableMigration, tables...)
}

// =============================================================================
// Health Check
// =============================================================================

// HealthCheck checks database connection health
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// Stats returns database connection pool statistics
func Stats(db *gorm.DB) (sql.DBStats, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return sql.DBStats{}, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Stats(), nil
}

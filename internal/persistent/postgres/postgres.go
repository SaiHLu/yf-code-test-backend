package postgres

import (
	"fmt"
	"log"
	"sync"
	"time"

	"codetest/internal/config"

	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	singletonDBInstance *gorm.DB
	mu                  sync.Mutex
)

type DBConnection struct {
	*gorm.DB
}

func NewPostgresDBConnection(cfg *config.AppConfig) (*DBConnection, error) {
	mu.Lock()
	defer mu.Unlock()

	if singletonDBInstance != nil {
		return &DBConnection{
			DB: singletonDBInstance,
		}, nil
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.POSTGRES_HOST,
		cfg.POSTGRES_PORT,
		cfg.POSTGRES_USERNAME,
		cfg.POSTGRES_PASSWORD,
		cfg.POSTGRES_DB,
		cfg.POSTGRES_SSLMODE)

	db, err := gorm.Open(pg.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// sqlDB.SetConnMaxIdleTime(time.Minute)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetMaxOpenConns(100)

	singletonDBInstance = db

	return &DBConnection{
		DB: singletonDBInstance,
	}, nil
}

func (db *DBConnection) Close() error {
	mu.Lock()
	defer mu.Unlock()

	if singletonDBInstance != nil {
		sqlDB, err := db.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func (db *DBConnection) LogConnectionStats() {
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Printf("Error getting DB stats: %v", err)
		return
	}

	// Get current stats
	stats := sqlDB.Stats()

	log.Printf("DB Connection Pool Stats:")
	log.Printf("- Open connections: %d", stats.OpenConnections)
	log.Printf("- In use connections: %d", stats.InUse)
	log.Printf("- Idle connections: %d", stats.Idle)
	log.Printf("- Wait count: %d", stats.WaitCount)
	log.Printf("- Wait duration: %v", stats.WaitDuration)
	log.Printf("- Max idle closed: %d", stats.MaxIdleClosed)
	log.Printf("- Max lifetime closed: %d", stats.MaxLifetimeClosed)
}

func (db *DBConnection) GetDBInstance() *gorm.DB {
	return singletonDBInstance
}

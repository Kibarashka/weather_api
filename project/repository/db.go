package repository

import (
	"fmt"
	"log"
	"time"
	"weather/project/config"
	"weather/project/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("repository.InitDB: failed to connect to database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("repository.InitDB: failed to get generic database object: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Println("Database connection established")
	return db, nil
}

func MigrateDB(db *gorm.DB) error {
	log.Println("Running database migrations...")
	err := db.AutoMigrate(
		&domain.Subscription{},
	)
	if err != nil {
		return fmt.Errorf("repository.MigrateDB: failed to run migrations: %w", err)
	}
	log.Println("Database migrations completed")
	return nil
}

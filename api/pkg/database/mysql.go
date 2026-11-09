package database

import (
	"api/config"
	"api/pkg/utils"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSQLConnection(cfg *config.Configs) (*gorm.DB, error) {
	// get DSN form ConnectionUrlBuilder
	dsn, err := utils.ConnectionUrlBuilder("mysql", cfg)
	if err != nil {
		return nil, fmt.Errorf("failded to Build MYSQL DNS: %w", err)
	}

	// logger Gorm base on Enviroment
	logMode := logger.Info
	if os.Getenv("APP_ENV") == "production" {
		logMode = logger.Warn
	}

	// Connect Databaes GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		return nil, fmt.Errorf("failded to Build MYSQL DNS: %w", err)
	}

	// get generic database object (sql.DB) for config Connection Pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.MySQL.ConnMaxLifetime)

	log.Println("MySQL Databases connected successfully!")
	return db, nil
}

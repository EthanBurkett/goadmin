package database

import (
	"time"

	"github.com/ethanburkett/goadmin/app/logger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data.sqlite"), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database:", zap.Error(err))
		return
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Error("Failed to get underlying DB:", zap.Error(err))
		return
	}

	// Set maximum number of open connections
	sqlDB.SetMaxOpenConns(25)

	// Set maximum number of idle connections
	sqlDB.SetMaxIdleConns(10)

	// Set maximum lifetime of a connection
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("Database connection pool configured",
		zap.Int("maxOpenConns", 25),
		zap.Int("maxIdleConns", 10),
		zap.Duration("connMaxLifetime", time.Hour),
	)
}

func AutoMigrate(models ...interface{}) {
	err := DB.AutoMigrate(
		models...,
	)
	if err != nil {
		logger.Error("Failed to auto migrate models:", zap.Error(err))
	}
}

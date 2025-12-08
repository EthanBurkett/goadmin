package database

import (
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
	}
}

func AutoMigrate(models ...interface{}) {
	err := DB.AutoMigrate(
		models...,
	)
	if err != nil {
		logger.Error("Failed to auto migrate models:", zap.Error(err))
	}
}

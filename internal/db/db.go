package db

import (
	"cogmoteHub/internal/models"
	"fmt"
	"log/slog"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func Init(host, user, password, dbName string) *gorm.DB {
	once.Do(func() {
		if host == "" || user == "" || dbName == "" {
			slog.Error("database connection details not provided")
			return
		}

		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
			host, user, password, dbName,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			slog.Error("failed to connect to database", "error", err)
			return
		}

		err = db.AutoMigrate(
			&models.Device{},
			&models.Animal{},
		)
		if err != nil {
			slog.Error("failed to auto-migrate database", "error", err)
			return
		}

		instance = db
		slog.Info("database connection initialized")
	})
	return instance
}

func Get() *gorm.DB {
	if instance == nil {
		slog.Warn("database connection requested before initialization")
	}
	return instance
}

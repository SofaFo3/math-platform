package database

import (
	"fmt"
	"log"
	"math-platform/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(host, port, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Автоматическая миграция (создание таблиц)
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	log.Println("✅ Database connected and migrated successfully")
	return db, nil
}

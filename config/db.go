package config

import (
	"errors"
	"log"

	"github.com/ahdaan98/go-gorm-crud/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg Config) {
	if cfg.DSN == "" {
		log.Fatal("Database url is empty")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	if err := Migration(DB); err != nil {
		log.Fatal(err)
	}
}

func Migration(DB *gorm.DB) error {
	if err := DB.AutoMigrate(&models.Book{}); err != nil {
		return errors.New("failed to create table")
	}
	return nil
}

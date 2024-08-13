package database

import (
	"log"

	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgres() *gorm.DB {
	dsn := config.Conf.Postgres.Dsn

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("new postgres: %v", err)
	}

	return db
}

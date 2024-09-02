package database

import (
	"time"

	"github.com/dwprz/prasorganic-user-service/src/common/log"
	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/sirupsen/logrus"
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
		log.Logger.WithFields(logrus.Fields{"location": "database.NewPostgres", "section": "gorm.Open"}).Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "database.NewPostgres", "section": "db.DB"}).Fatal(err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db
}

func ClosePostgres(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "database.ClosePostgres", "section": "db.DB"}).Error(err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "database.ClosePostgres", "section": "sqlDB.Close"}).Error(err)
	}
}

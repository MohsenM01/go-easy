package database

import (
	"go-easy/config"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	return db
}

func InitDb() {
	once.Do(func() {

		cfg := config.LoadConfig()

		dsn := cfg.DbConnection

		gormConfig := &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		}

		var err error
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
		if err != nil {
			log.Fatalf("Can not connect to DataBase : %v", err)
		}

		migrations(db)

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Can not connect to DataBase: %v", err)
		}

		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
		sqlDB.SetConnMaxIdleTime(2 * time.Minute)

		log.Println("Connect to database was successfull")
	})
}

func migrations(db *gorm.DB) {
}

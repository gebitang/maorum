package models

import (
	"context"
	"gebitang.com/maorum/models/dbm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var (
	Store          *gorm.DB
	DailyItemStore dbm.DailyItemStore
)

func init() {
	// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open("rum.db"), &gorm.Config{})
	if err != nil {
		log.Println("err ", err)
		os.Exit(1)
	}
	sqlDB, err := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	db.AutoMigrate(&dbm.DailyItem{})
	Store = db
	DailyItemStore = &dailyItemStore{}
}

func GetDB(ctx context.Context) *gorm.DB {
	return Store.WithContext(ctx)
}

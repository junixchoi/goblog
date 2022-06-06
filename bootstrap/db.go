package bootstrap

import (
	"goblog/pkg/model"
	"time"
)

func SetupDB() {
	db := model.ConnectDB()

	sqlDB, _ := db.DB()

	sqlDB.SetMaxOpenConns(100)

	sqlDB.SetMaxIdleConn(25)

	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}

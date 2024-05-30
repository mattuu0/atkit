package database

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	isInit = false
	conn *gorm.DB = nil
)

func Init() {
	//データベース接続
	dbconn,err := gorm.Open(sqlite.Open(os.Getenv("test.db")),&gorm.Config{})

	//エラー処理
	if err != nil {
		log.Fatalf("failed to connect : %s",err)
	}

	//マイグレーション
	dbconn.AutoMigrate(
		&User{},
	)

	//グローバル変数に格納
	conn = dbconn

	isInit = true
}
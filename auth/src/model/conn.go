package model

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
	dbconn,err := gorm.Open(sqlite.Open(os.Getenv("DB_PATH")),&gorm.Config{})

	//エラー処理
	if err != nil {
		log.Fatalf("failed to connect : %s",err)
	}

	//マイグレーション
	dbconn.AutoMigrate(
		&User{},
		&Session{},
	)

	//グローバル変数に格納
	conn = dbconn

	isInit = true
}

func GetConn() *gorm.DB {
	//初期化
	if !isInit {
		Init()
	}

	return conn
}
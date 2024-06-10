package main

import (
	"auth/auth"
	"auth/controller"
	"auth/model"
	"auth/service"
	"auth/session"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// .envを呼び出します。
func loadEnv() {
	// ここで.envファイル全体を読み込みます。
	// この読み込み処理がないと、個々の環境変数が取得出来ません。
	// 読み込めなかったら err にエラーが入ります。
	err := godotenv.Load(".env")
	
	// もし err がnilではないなら、"読み込み出来ませんでした"が出力されます。
	if err != nil {
		log.Fatalf("読み込み出来ませんでした: %v", err)
	} 
}

func Init() {
	loadEnv()

	// データベース接続
	model.Init()

	// セッション初期化
	session.Init(os.Getenv("JWT_SECRET"))

	//JWt 初期化
	auth.Init()


	//Oauth初期化
	controller.Oauth_Init()

	service.Init()
}
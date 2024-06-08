package main

import (
	authsdk "atkit/auth_sdk"
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
	//env よみこみ　
	loadEnv()

	//認証SDK初期化
	err := authsdk.Init(os.Getenv("PUB_KEY_PATH"))

	//エラー処理
	if err != nil {
		log.Println(err)
		log.Fatalln("cannot init auth sdk")
	}
}
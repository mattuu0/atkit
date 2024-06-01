package main

import (
	"auth/oauth"
	"auth/session"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	//初期化
	Init()

	//session.GetSession("aiueo","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.142.86 Safari/537.36","127.0.0.1")
	session.GenToken("aiueo")
}

func ServerMain() {

	//ルーター設定
	router := gin.Default()

	//Oauth 設定
	oauth.Oauth_Setup(router)

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run(os.Getenv("BindAddr")) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	//初期化
	Init()

	//ルーター設定
	router := gin.Default()

	//Oauth 設定
	Oauth_Setup(router)

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run(os.Getenv("BindAddr")) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}




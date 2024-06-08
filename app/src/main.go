package main

import (
	authsdk "atkit/auth_sdk"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//初期化
	Init()

	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/authed", authsdk.AuthMiddleware(), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "auth success",
		})
	})

	router.Run(":3001") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

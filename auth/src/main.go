package main

import (
	"auth/route"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	Init()

	//ルーター設定
	router := route.GenRouter()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run(os.Getenv("BindAddr")) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
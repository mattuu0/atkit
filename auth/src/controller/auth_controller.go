package controller

import (
	"auth/model"
	"auth/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GenAccessToken(ctx *gin.Context) {
	//ユーザー取得
	user,_ := ctx.Get("user")

	//トークン作成
	atoken,err := service.GenToken(user.(model.User).UserID)

	//エラー処理
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": atoken,
	})
}
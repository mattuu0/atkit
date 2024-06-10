package controller

import (
	"auth/model"
	"auth/service"
	"auth/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(ctx *gin.Context) {
	//セッション取得
	now_session,_ := ctx.Get("session")

	//ログアウト
	err := service.Logout(now_session.(*model.Session))

	//エラー処理
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}


func StartUpdateSession(ctx *gin.Context) {
	//セッション取得
	now_session,_ := ctx.Get("session")

	//セッション更新
	update_token,err := service.StartUpdateSession(now_session.(*model.Session),ctx.Request.UserAgent(),ctx.ClientIP())

	//エラー処理
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//cookie 上書き
	ctx.SetCookie("session_token",update_token,int(util.GetExp()),"/","localhost",true,true)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func SubmitUpdateSession(ctx *gin.Context) {
	//セッション取得
	now_session,_ := ctx.Get("session")

	//セッション更新
	err := service.SubmitUpdateSession(now_session.(*model.Session),ctx.Request.UserAgent(),ctx.ClientIP())

	//エラー処理
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
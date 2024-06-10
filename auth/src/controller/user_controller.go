package controller

import (
	"auth/model"
	"auth/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//ユーザー取得
func GetUser(ctx *gin.Context) {
	//ユーザー取得
	ctx.JSON(http.StatusOK, ctx.MustGet("user").(model.User))
}


//アイコン取得
func GetIcon(ctx *gin.Context) {
	//ユーザー取得
	user,err := service.GetUser(ctx.Param("userid"))

	//エラー処理
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//画像を読み込む
	ctx.File(user.IconPath)
}

//アイコンアップロード
func UploadIcon(ctx *gin.Context) {
	//ユーザー取得
	user,_ := ctx.MustGet("user").(model.User)

	//アイコン取得
	icon_file, err := ctx.FormFile("icon")
	if err != nil {
		log.Println(err)

		ctx.HTML(http.StatusBadRequest, "oauth/oauth_error.html", gin.H{
			"error_log": err.Error(),
		})
		return
	}

	//アイコンをアップロード
	err = service.UploadIcon(user.UserID,icon_file)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.HTML(http.StatusBadRequest, "oauth/oauth_error.html", gin.H{
			"error_log": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
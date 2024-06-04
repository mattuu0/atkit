package main

import (
	"auth/database"
	"auth/oauth"
	"auth/session"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	Init()

	ServerMain()
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

	//グループ作成
	authed_group := router.Group("/authed")
	{
		//認証済みAPI
		authed_group.Use(Middleware())

		authed_group.GET("/GetUser", func(ctx *gin.Context) {
			//ユーザー取得
			ctx.JSON(http.StatusOK, ctx.MustGet("user").(database.User))
		})

		//ログアウト
		authed_group.POST("/Logout", func(ctx *gin.Context) {
			//セッション取得
			now_session,_ := ctx.Get("session")

			//セッション削除
			err := session.DeleteSession(now_session.(*database.Session).SessionID)

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
		})

		authed_group.POST("/Update", func(ctx *gin.Context) {
			//セッション取得
			now_session,_ := ctx.Get("session")

			//セッション更新
			update_token,err := session.UpdateSession(now_session.(*database.Session).TokenID,ctx.Request.UserAgent(),ctx.ClientIP())

			//エラー処理
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			//cookie 上書き
			ctx.SetCookie("session_token",update_token,int(session.GetExp()),"/","localhost",true,true)

			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
			})
		})

		authed_group.POST("/SubmitUpdate", func(ctx *gin.Context) {
			//セッション取得
			now_session,_ := ctx.Get("session")

			//セッション更新
			err := session.SubmitUpdate(now_session.(*database.Session).TokenID,ctx.Request.UserAgent(),ctx.ClientIP())

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
		})
	}

	router.Run(os.Getenv("BindAddr")) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

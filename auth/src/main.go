package main

import (
	"auth/auth"
	"auth/database"
	"auth/oauth"
	"auth/session"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	nocache "github.com/alexander-melentyev/gin-nocache"
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

		//更新を開始するエンドポイント	
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

		//更新を確定するエンドポイント
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

		//アクセストークン作成
		authed_group.GET("/GenToken",nocache.NoCache(), func(ctx *gin.Context) {
			//ユーザー取得
			user,_ := ctx.Get("user")

			//トークン作成
			atoken,err := auth.GenAccessToken(user.(database.User).UserID)

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
		})
	}

	//アイコングループ
	uicon_group := router.Group("/uicon")
	{
		//アイコンを取得するエンドポイント
		uicon_group.GET("/:userid",func(ctx *gin.Context) {
			//ユーザー取得
			user,err := database.GetUserByID(ctx.Param("userid"))
	
			//エラー処理
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
	
			//画像を読み込む
			ctx.File(user.IconPath)
		})

		//アイコンを更新するエンドポイント
		uicon_group.POST("/upicon",Middleware(),func(ctx *gin.Context) {
			//ユーザー取得
			user,_ := ctx.MustGet("user").(database.User)

			//アイコン取得
			icon_file, err := ctx.FormFile("icon")
			if err != nil {
				log.Println(err)

				ctx.HTML(http.StatusBadRequest, "oauth/oauth_error.html", gin.H{
					"error_log": err.Error(),
				})
				return
			}

			//ファイルパス指定
			savepath := filepath.Join("./assets/UserIcons",user.UserID + ".png")

			//アイコンファイルを開く
			ofile,err := icon_file.Open()

			//エラー処理
			if err != nil {
				log.Println(err)

				ctx.HTML(http.StatusBadRequest, "oauth/oauth_error.html", gin.H{
					"error_log": err.Error(),
				})
				return
			}

			//アイコンファイルを保存
			err = oauth.Resizeio(ofile,savepath)

			//エラー処理
			if err != nil {
				log.Println(err)

				ctx.HTML(http.StatusBadRequest, "oauth/oauth_error.html", gin.H{
					"error_log": err.Error(),
				})
				return
			}

			//データーベース
			dbconn := database.GetConn()

			//ユーザーを取得
			user.IconPath = savepath
			
			//ユーザー更新
			result := dbconn.Save(&user)

			//エラー処理
			if result.Error != nil {
				log.Println(result.Error)
				ctx.HTML(http.StatusBadRequest, "oauth/oauth_error.html", gin.H{
					"error_log": result.Error.Error(),
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

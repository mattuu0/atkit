package route

import (
	"auth/controller"
	"auth/middleware"

	nocache "github.com/alexander-melentyev/gin-nocache"
	"github.com/gin-gonic/gin"
)

func GenRouter() *gin.Engine {
	router := gin.Default()

	//HTMLを読み込む
	router.LoadHTMLFiles("./htmls/oauth/oauth_error.html")

	//oauth
	oauth_group := router.Group("/oauth")
	{
		//Oauth認証
		oauth_group.GET("/:provider", controller.Oauth)
		oauth_group.GET("/:provider/callback", controller.Oauth_Callback)
	}

	//グループ作成
	authed_group := router.Group("/authed")
	{
		//認証済みAPI
		authed_group.Use(middleware.Auth_Middleware())

		authed_group.GET("/GetUser",controller.GetUser)

		//ログアウト
		authed_group.POST("/Logout", controller.Logout)

		//更新を開始するエンドポイント	
		authed_group.POST("/Update", controller.StartUpdateSession)

		//更新を確定するエンドポイント
		authed_group.POST("/SubmitUpdate", controller.SubmitUpdateSession)

		//アクセストークン作成
		authed_group.GET("/GenToken",nocache.NoCache(), controller.GenAccessToken)
	}


	//アイコングループ
	uicon_group := router.Group("/uicon")
	{
		//アイコンを取得するエンドポイント
		uicon_group.GET("/:userid",controller.GetIcon)

		//アイコンを更新するエンドポイント
		uicon_group.POST("/upicon",middleware.Auth_Middleware(),controller.UploadIcon)
	}

	return router
}
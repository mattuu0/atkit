package main

import (
	"context"
	"log"
	"os"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth"

	//" github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/discord"
)

func Oauth_Init() {
	goth.UseProviders(
		//google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"),os.Getenv("GOOGLE_CALLBACK_URL")),
		discord.New(os.Getenv("DISCORD_KEY"), os.Getenv("DISCORD_SECRET"),os.Getenv("DISCORD_CALLBACK_URL"),"identify","email"),
	)
}

func Oauth_Setup(router *gin.Engine) {
	//Oauth 初期化
	Oauth_Init()

	//HTMLを読み込む
	router.LoadHTMLFiles("./htmls/oauth/oauth_error.html")

	//グループを作成
	ogroup := router.Group("/oauth")
	{
		//認証を開始するエンドポイント
		ogroup.GET("/:provider", func(ctx *gin.Context) {
			provider := ctx.Param("provider")
			ctx.Request = contextWithProviderName(ctx, provider)
	
			gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
		})

		//認証を完了するエンドポイント
		ogroup.GET("/:provider/callback", func(ctx *gin.Context) {
			provider := ctx.Param("provider")
			ctx.Request = contextWithProviderName(ctx, provider)
	
			//認証を完了する
			user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)

			//エラー処理
			if err != nil {
				log.Println(ctx.Writer, err)
				ctx.HTML(http.StatusInternalServerError,"oauth_error.html",gin.H{
					"error_log" : err.Error(),
				})
				return
			}

			//ユーザを表示する
			log.Printf("%#v", user)
		})
	}
}

//プロバイダを取得する
func contextWithProviderName(ctx *gin.Context, provider string) (*http.Request){
    return  ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", provider))
}

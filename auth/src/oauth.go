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
		discord.New(os.Getenv("DISCORD_KEY"), os.Getenv("DISCORD_SECRET"),os.Getenv("DISCORD_CALLBACK_URL")),
	)
}

func Oauth_Setup(router *gin.Engine) {
	//Oauth 初期化
	Oauth_Init()

	ogroup := router.Group("/oauth")
	{
		ogroup.GET("/:provider", func(ctx *gin.Context) {
			provider := ctx.Param("provider")
			ctx.Request = contextWithProviderName(ctx, provider)
	
			gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
		})

		ogroup.GET("/:provider/callback", func(ctx *gin.Context) {
			provider := ctx.Param("provider")
			ctx.Request = contextWithProviderName(ctx, provider)
	
			user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
			if err != nil {
				log.Println(ctx.Writer, err)
				return
			}

			log.Printf("%#v", user)
		})
	}
}


func contextWithProviderName(ctx *gin.Context, provider string) (*http.Request){
    return  ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", provider))
}

package controller

import (
	"auth/service"
	"auth/util"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/microsoftonline"
)

func Oauth_Init() {
	log.Println("Oauth_Init")
	
	goth.UseProviders(
		google.New(os.Getenv("Google_KEY"), os.Getenv("Google_SECRET"), os.Getenv("Google_CALLBACK_URL"), "email", "profile"),
		discord.New(os.Getenv("DISCORD_KEY"), os.Getenv("DISCORD_SECRET"), os.Getenv("DISCORD_CALLBACK_URL"), "identify", "email"),
		microsoftonline.New(os.Getenv("Microsoft_KEY"), os.Getenv("Microsoft_SECRET"), os.Getenv("Microsoft_CALLBACK_URL"), "openid", "profile", "email"),
	)
}

//Oauth 認証 開始
func Oauth(ctx *gin.Context) {
	provider := ctx.Param("provider")
	ctx.Request = contextWithProviderName(ctx, provider)

	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func Oauth_Callback(ctx *gin.Context) {
	//プロバイダ取得
	provider := ctx.Param("provider")
	ctx.Request = contextWithProviderName(ctx, provider)

	//認証を完了する
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)

	//エラー処理
	if err != nil {
		log.Println(ctx.Writer, err)
		ctx.HTML(http.StatusInternalServerError, "oauth_error.html", gin.H{
			"error_log": err.Error(),
		})
		return
	}

	//サービスに引き渡す
	session_token,err := service.Callback_Oauth(user, ctx.Request.UserAgent(), ctx.ClientIP())
	
	//エラー処理
	if err != nil {
		log.Println(ctx.Writer, err)
		ctx.HTML(http.StatusInternalServerError, "oauth_error.html", gin.H{
			"error_log": err.Error(),
		})
		return
	}

	//Cookie に設定
	ctx.SetCookie("session_token", session_token, int(util.GetExp()), "/", "localhost", true, true)

	ctx.Redirect(http.StatusFound, "/statics/index.html")
}

// プロバイダを取得する
func contextWithProviderName(ctx *gin.Context, provider string) *http.Request {
	return ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", provider))
}

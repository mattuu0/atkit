package oauth

import (
	"auth/database"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/google"
)

func Oauth_Init() {
	goth.UseProviders(
		google.New(os.Getenv("Google_KEY"), os.Getenv("Google_SECRET"),os.Getenv("Google_CALLBACK_URL"),"email","profile"),
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
			//プロバイダ取得
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

			database.CreateUser(database.User{
				UserID: user.UserID,
				Provider: provider,
				UserName: user.Name,
				NickName: user.NickName,
				Email: user.Email,
				IconPath: user.AvatarURL,
				ProviderUID: user.UserID,
				IsVeriry: false,
			})

			//ユーザを表示する
			log.Printf("%#v", user)
		})
	}
}

//プロバイダを取得する
func contextWithProviderName(ctx *gin.Context, provider string) (*http.Request){
    return  ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", provider))
}

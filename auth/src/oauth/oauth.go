package oauth

import (
	"auth/database"
	"auth/session"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"

	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/microsoftonline"
)

func Oauth_Init() {
	goth.UseProviders(
		google.New(os.Getenv("Google_KEY"), os.Getenv("Google_SECRET"), os.Getenv("Google_CALLBACK_URL"), "email", "profile"),
		discord.New(os.Getenv("DISCORD_KEY"), os.Getenv("DISCORD_SECRET"), os.Getenv("DISCORD_CALLBACK_URL"), "identify", "email"),
		microsoftonline.New(os.Getenv("Microsoft_KEY"), os.Getenv("Microsoft_SECRET"), os.Getenv("Microsoft_CALLBACK_URL"), "openid", "profile", "email"),
	)

	//セッション初期化
	session.Init()
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

			//データベース
			dbconn := database.GetConn()

			//エラー処理
			if err != nil {
				log.Println(ctx.Writer, err)
				ctx.HTML(http.StatusInternalServerError, "oauth_error.html", gin.H{
					"error_log": err.Error(),
				})
				return
			}

			//ユーザーを取得する
			usr, err := database.GetUser(user.Provider, user.UserID)

			if err == gorm.ErrRecordNotFound {
				//見つからないときユーザーを作成する
				err := database.CreateUser(database.User{
					UserID:      user.UserID,
					Provider:    user.Provider,
					UserName:    user.Name,
					NickName:    user.NickName,
					Email:       user.Email,
					IconPath:    "./assets/UserIcons/default.png",
					ProviderUID: user.UserID,
					IsVeriry:    false,
				})

				//エラー処理
				if err != nil {
					log.Println(err)
					ctx.HTML(http.StatusInternalServerError, "oauth_error.html", gin.H{
						"error_log": err.Error(),
					})
					return
				}

				//ユーザーを取得
				get_usr, err := database.GetUser(user.Provider, user.UserID)

				//エラー処理
				if err != nil {
					log.Println(err)
					ctx.HTML(http.StatusInternalServerError, "oauth_error.html", gin.H{
						"error_log": err.Error(),
					})
					return
				}

				//アイコンパス
				//アイコンのURLが存在するか
				if user.AvatarURL != "" {
					iconpath := "./assets/UserIcons/" + get_usr.UserID + ".png"
					//アイコンを保存する
					err := SaveIcon(user.AvatarURL, iconpath)

					//エラー処理
					if err != nil {
						log.Println(ctx.Writer, err)
						ctx.HTML(http.StatusInternalServerError, "oauth_error.html", gin.H{
							"error_log": err.Error(),
						})
						return
					}

					//ユーザーを更新する
					get_usr.IconPath = iconpath

					//ユーザーを更新
					result := dbconn.Save(&get_usr)

					//エラー処理
					if result.Error != nil {
						log.Println(ctx.Writer, result.Error)
						ctx.HTML(http.StatusInternalServerError, "oauth_error.html", gin.H{
							"error_log": result.Error.Error(),
						})
						return
					}
				}

				//ユーザーを取得
				usr = get_usr
			} else if err != nil {
				//エラー処理
				ctx.HTML(http.StatusInternalServerError, "oauth_error.html", gin.H{
					"error_log": err.Error(),
				})
				return
			}


			//アイコンファイルが存在するか
			if _, err := os.Stat(usr.IconPath); err != nil {
				//存在しないとき
				//アイコンURLが存在するか
				if user.AvatarURL != "" {
					//アイコンを保存する
					err := SaveIcon(user.AvatarURL, usr.IconPath)

					//エラー処理
					if err != nil {
						log.Println(err)
						ctx.HTML(http.StatusInternalServerError, "oauth_error.html", gin.H{
							"error_log": err.Error(),
						})
						return
					}
				}
			}

			//セッションを作成する
			session_token, err := session.GetSession(usr.UserID, ctx.Request.UserAgent(), ctx.ClientIP())

			//エラー処理
			if err != nil {
				log.Println(ctx.Writer, err)
				ctx.HTML(http.StatusInternalServerError, "oauth_error.html", gin.H{
					"error_log": err.Error(),
				})
				return
			}

			//Cookie に設定
			ctx.SetCookie("session_token", session_token, int(session.GetExp()), "/", "localhost", true, true)

			ctx.Redirect(http.StatusFound, "/statics/index.html")
		})
	}
}

// プロバイダを取得する
func contextWithProviderName(ctx *gin.Context, provider string) *http.Request {
	return ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", provider))
}

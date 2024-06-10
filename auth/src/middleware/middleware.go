package middleware

import (
	"auth/model"
	"auth/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth_Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("success",false)
		ctx.Set("user",model.User{})
		ctx.Set("session",model.Session{})

		//トークン取得
		token,err := GetToken(ctx)

		//エラー処理
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//トークン検証
		if token == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//トークン検証
		tokenid,err := session.VerifyToken(token)

		//エラー処理
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//セッションを取得する
		get_session,err := model.GetSessionByTokenID(tokenid)

		//エラー処理
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//ユーザーを取得する
		get_user,err := model.GetUserByID(get_session.UserID)

		//エラー処理
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//状態を設定する
		ctx.Set("success",true)
		ctx.Set("user",*get_user)
		ctx.Set("session",get_session)

		ctx.Next()
	}
}

func GetToken(ctx *gin.Context) (string,error) {
	//トークン取得
	return ctx.Cookie("session_token")
}
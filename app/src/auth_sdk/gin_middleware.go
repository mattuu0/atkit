package authsdk

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//初期化
		ctx.Set("userid", "")

		//Bearerトークンを取得
		bearer_token := ctx.Request.Header.Get("Authorization")
		if bearer_token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token"})
			ctx.Abort()
			return
		}

		//長さ取得
		if len(bearer_token) < 7 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		//Bearer 確認
		if bearer_token[:7] != "Bearer " {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		//トークンを取得
		token := bearer_token[7:]

		//トークンを検証
		userid, err := VerifyAccessToken(token)

		//エラー処理
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		//ユーザーIDを格納
		ctx.Set("userid", userid)
		ctx.Next()
	}
}
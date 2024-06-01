package session

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//トークン生成
func GenToken(tokenid string) {
	// JWTに付与する構造体
	claims := jwt.MapClaims{
		"tokenid": tokenid,
		"exp": GetExp(), // 72時間が有効期限
	}

	// ヘッダーとペイロード生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// トークンに署名を付与
	accessToken, _ := token.SignedString([]byte("ACCESS_SECRET_KEY"))
	log.Println("accessToken:", accessToken)
}

func GetExp() int64 {
	return time.Now().AddDate(1,0,0).Unix()
}
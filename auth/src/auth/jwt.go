package auth

import (
	"crypto/ed25519"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	//JWT 秘密鍵
	priv_key ed25519.PrivateKey = nil

	//アクセストークン有効期限 (5分)
	AccessExp int64 = 300
)

func Init() {
	//秘密鍵と公開鍵を生成
	gen_priv, err := GenKey(os.Getenv("Ed25519_KEY_PATH"), os.Getenv("Ed25519_PUB_PATH"))

	//エラー処理
	if err != nil {
		log.Println(err)
		return
	}

	//グローバル変数に格納
	priv_key = gen_priv
}

func GenAccessToken(userid string) (string, error) {
	//JWTを生成
	claims := jwt.MapClaims{
		"userid" : userid,
		"exp" : time.Now().Unix() + AccessExp,
	}

	//JWTを生成
	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, claims)

	//JWTを署名
	signed_token, err := token.SignedString(priv_key)

	//エラー処理
	if err != nil {
		return "", err
	}

	return signed_token, nil
}

package authsdk

import (
	"crypto/ed25519"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var (
	pub_key ed25519.PublicKey = nil
)

func Init(pubkey_path string) error {
	//鍵を読み込む
	read_key,err := ReadKey(pubkey_path)

	//エラー処理
	if err != nil {
		return err
	}

	//グローバル変数に格納
	pub_key = read_key

	return nil
}

//トークンを検証 ユーザーIDを返す
func VerifyAccessToken(tokenString string) (string, error) {
	//JWTを検証
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pub_key, nil
	})

	//JWTを検証
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["userid"].(string), nil
	} else {
		return "", err
	}
}
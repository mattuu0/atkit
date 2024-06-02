package session

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	//TODO JWT シークレット 実際に使う際は変える
	secret = "VHz4OnFGuftXGx9FDOA9AMyORaclIGbqi1W8dT6bql7BOkprjRZa4KCay8DqNRYE"
)

// トークン生成
func GenToken(tokenid string) (string,error) {
	// JWTに付与する構造体
	claims := jwt.MapClaims{
		"tokenid": tokenid,
		"exp":     GetExp(), //有効期限
	}

	// ヘッダーとペイロード生成
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	// トークンに署名を付与
	signed_token, err := token.SignedString(priv_key)

	// エラー処理
	if err != nil {
		return "",err
	}

	return signed_token,nil
}

func VerifyToken(tokenString string) (string,error) {
	//トークンを検証する
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		//検証出来たら鍵を返す
		return pub_key, nil
	})

	//トークンが認証されているか確認
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//トークンIDを返す
		return claims["tokenid"].(string),nil
	} else {
		//エラーを返す
		return "",err
	}
}

func GetExp() int64 {
	return time.Now().AddDate(1, 0, 0).Unix()
}

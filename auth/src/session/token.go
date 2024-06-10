package session

import (
	"auth/util"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// トークン生成
func GenToken(tokenid string) (string,error) {
	// JWTに付与する構造体
	claims := jwt.MapClaims{
		"tokenid": tokenid,
		"exp":     util.GetExp(), //有効期限
	}

	// ヘッダーとペイロード生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// トークンに署名を付与
	signed_token, err := token.SignedString([]byte(JwtSecret))

	// エラー処理
	if err != nil {
		return "",err
	}

	return signed_token,nil
}

func VerifyToken(tokenString string) (string,error) {
	//トークンを検証する
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		//検証出来たら鍵を返す
		return []byte(JwtSecret), nil
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

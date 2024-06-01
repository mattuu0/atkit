package session

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	//TODO JWT シークレット 実際に使う際は変える
	secret = "VHz4OnFGuftXGx9FDOA9AMyORaclIGbqi1W8dT6bql7BOkprjRZa4KCay8DqNRYE"
)

// トークン生成
func GenToken(tokenid string) {
	// JWTに付与する構造体
	claims := jwt.MapClaims{
		"tokenid": tokenid,
		"exp":     GetExp(), // 72時間が有効期限
	}

	// ヘッダーとペイロード生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// トークンに署名を付与
	signed_token, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("accessToken:", signed_token)

	VerifyToken(signed_token,)
}

func VerifyToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Println(claims["tokenid"])
	} else {
		fmt.Println(err)
	}

	return false
}

func GetExp() int64 {
	return time.Now().AddDate(1, 0, 0).Unix()
}

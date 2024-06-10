package service

import "auth/auth"

// トークン生成
func GenToken(userid string) (string, error) {
	return auth.GenAccessToken(userid)
}
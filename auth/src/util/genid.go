package util

import "github.com/google/uuid"

func GenID() string {
	//UUID生成
	uid,err := uuid.NewRandom()

	//エラー処理
	if err != nil {
		return ""
	}

	return uid.String()
}
package util

import "time"


func GetExp() int64 {
	return time.Now().AddDate(1, 0, 0).Unix()
}

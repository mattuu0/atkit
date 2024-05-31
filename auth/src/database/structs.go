package database

type Session struct {
	//セッションID
	SessionID string `gorm:"primaryKey"`

	//ユーザーID
	UserID string

	//トークンID
	TokenID string

	//UserAgent
	UserAgent string

	//IPアドレス
	IPAddress string
}
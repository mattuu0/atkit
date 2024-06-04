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

	//更新中か否
	IsUpdating bool

	//アップデートID
	UpdateID string `gorm:"default:null"`

	//タイプ
	Type string	`gorm:"default:access"`

	//有効期限
	Exp int64 `gorm:"default:-1"`	
}
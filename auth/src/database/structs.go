package database

type User struct {
	//ユーザーID
	UserID string `gorm:"primaryKey"`

	//認証に使ったプロバイダ
	Provider string

	//ユーザー名
	UserName string

	//ニックネーム
	NickName string

	//メールアドレス
	Email string

	//アイコンパス
	IconPath string `gorm:"default:default.png"`

	//プロバイダのユーザーID
	ProviderUID string

	//メールアドレスが検証されているか
	IsVeriry bool `gorm:"default:false"`
}

package database

import (
	"github.com/google/uuid"
)

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

func CreateUser(user User) error {
	//データベース接続
	dbconn := GetConn()

	//ユーザーID生成
	uid,err := uuid.NewRandom()

	//エラー処理
	if err != nil {
		return err
	}

	//ユーザーID生成
	user.UserID = uid.String()

	//ユーザー作成
	result := dbconn.Create(&user)

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetUser(provider, uid string) (*User,error) {
	//データベース接続
	dbconn := GetConn()

	var user User

	//ユーザー取得
	result := dbconn.Where(User{
		Provider: provider,
		ProviderUID: uid,
	}).First(&user)

	//エラー処理
	if result.Error != nil {
		return nil,result.Error
	}

	return &user,nil
}

func GetUserByID(uid string) (*User,error) {
	//データベース接続
	dbconn := GetConn()

	var user User

	//ユーザー取得
	result := dbconn.Where(User{
		UserID: uid,
	}).First(&user)

	//エラー処理
	if result.Error != nil {
		return nil,result.Error
	}

	return &user,nil
}
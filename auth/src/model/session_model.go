package model

import "errors"

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

//有効期限切れのセッションを取得する
func GetExperiedSession(exp int64) ([]Session, error) {
	//データベース接続
	dbconn := GetConn()
		
	//結果を格納する変数
	sessions := []Session{}

	//db検索
	result := dbconn.Where("exp < ?", exp).Find(&sessions)

	//エラー処理
	if result.Error != nil {
		return []Session{}, result.Error
	}

	return sessions, nil
}


// IDからセッション取得
func GetSessionByTokenID(tokenid string) (*Session, error) {
	//データベース接続
	dbconn := GetConn()

	//セッション取得
	var get_session Session

	//セッション取得
	result := dbconn.Where(Session{
		TokenID: tokenid,
	}).First(&get_session)

	//エラー処理
	if result.Error != nil {
		return nil, result.Error
	}

	return &get_session, nil
}

// IDでセッション削除
func DeleteSession(sessionid string) error {
	//データベース接続
	dbconn := GetConn()

	//セッション削除
	result := dbconn.Where(Session{
		SessionID: sessionid,
	}).Unscoped().Delete(&Session{})

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	return nil
}

//セッションを上書きする
func UpdateSession(session *Session) error {
	//セッションがnilならエラー
	if session == nil {
		return errors.New("session is nil")
	}

	//データベース接続
	dbconn := GetConn()

	//セッション保存
	result := dbconn.Save(session)

	return result.Error
}

//セッションを作成する
func CreateSession(session *Session) error {
	//セッションがnilならエラー
	if session == nil {
		return errors.New("session is nil")
	}

	//データベース接続
	dbconn := GetConn()

	//セッション作成
	result := dbconn.Create(session)

	return result.Error
}
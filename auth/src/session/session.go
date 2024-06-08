package session

import (
	"auth/database"
	"errors"
	"log"
	"os"
	"time"
)

var (
	//JWT 鍵
	JwtSecret string = ""
	//無期限
	noexp int64 = -1

	//更新用トークン有効期限 (5分)
	refresh_exp int64 = 300
)

func Init() {
	//データベース接続
	database.Init()

	//JWT鍵の読み込み
	JwtSecret = os.Getenv("JWT_SECRET")

	//有効期限切れセッションを削除する関数
	go func() {
		defer recover()

		//セッション削除
		Remove_Expired_Session()
	}()
}

func Remove_Expired_Session() {
	//データベース接続
	dbconn := database.GetConn()

	//古いセッション取得
	for {
		//3秒待機
		time.Sleep(time.Second * 3)

		sessions := []database.Session{}

		//db検索
		result := dbconn.Where("exp < ?", time.Now().Unix()).Find(&sessions)

		//エラー処理
		if result.Error != nil {
			log.Println(result.Error)
			continue
		}

		//古いセッションを削除
		for _, session := range sessions {
			log.Println(session.Exp)

			//有効期限が設定されていないとき
			if session.Exp == noexp {
				continue
			}

			//アクセストークンの場合
			if session.Type == "access" {
				//セッション取得
				get_session, err := GetSessionByTokenID(session.UpdateID)

				//エラー処理
				if err != nil {
					log.Println(err)
					continue
				}

				//セッション削除
				result = dbconn.Unscoped().Delete(&get_session)

				//エラー処理
				if result.Error != nil {
					log.Println(result.Error)
					continue
				}

				continue
			}

			//更新元セッション取得
			old_session, err := GetSessionByTokenID(session.UpdateID)

			//エラー処理
			if err == nil {
				log.Println(err)

				//更新中を外す
				old_session.IsUpdating = false

				//更新
				result := dbconn.Save(&old_session)

				//エラー処理
				if result.Error != nil {
					log.Println(result.Error)
					continue
				}
			} else if err != nil {
				log.Println(err)
			}

			//セッション削除
			result = dbconn.Unscoped().Delete(&session)

			//エラー処理
			if result.Error != nil {
				log.Println(result.Error)
				continue
			}
		}

	}
}

// セッションを作成する (トークンを返す)
func GetSession(bindid string, useragent string, ipaddr string) (string, error) {
	//データベース接続
	dbconn := database.GetConn()

	//トークンID生成
	tokenid := GenID()

	//トークン取得
	stoken, err := GenToken(tokenid)

	//エラー処理
	if err != nil {
		return "", err
	}

	//セッションID取得
	SessionID := GenID()

	//セッション作成
	session_data := database.Session{
		SessionID:  SessionID,
		UserID:     bindid,
		TokenID:    tokenid,
		UserAgent:  useragent,
		IPAddress:  ipaddr,
		IsUpdating: false,
		Exp:        GetExp(),
	}

	//セッション作成
	result := dbconn.Create(&session_data)

	//エラー処理
	if result.Error != nil {
		return "", result.Error
	}

	return stoken, nil
}

// 更新用のトークンを返す
// セッションを更新する
func UpdateSession(tokenid string, useragent string, ipaddr string) (string, error) {
	//データベース接続
	//dbconn := database.GetConn()

	//セッション取得
	session, err := GetSessionByTokenID(tokenid)

	//エラー処理
	if err != nil {
		return "", err
	}

	//更新用のセッションの場合
	if session.Type == "refresh" {
		return "", errors.New("session is refresh")
	}

	//更新中の場合
	if session.IsUpdating {
		return "", errors.New("session is updating")
	}

	//トークンID生成
	new_tokenid := GenID()
	//更新用トークン生成
	new_token, err := GenToken(new_tokenid)

	//エラー処理
	if err != nil {
		return "", err
	}

	//更新用セッション作成
	//データベース接続
	dbconn := database.GetConn()

	//セッションID取得
	SessionID := GenID()

	//セッション作成 (有効期限 5分)
	session_data := database.Session{
		SessionID:  SessionID,
		UserID:     session.UserID,
		TokenID:    new_tokenid,
		UserAgent:  useragent,
		IPAddress:  ipaddr,
		IsUpdating: false,
		UpdateID:   tokenid,
		Type:       "refresh",
		Exp:        time.Now().Add(time.Second * time.Duration(refresh_exp)).Unix(),
	}

	//セッション作成
	result := dbconn.Create(&session_data)

	//エラー処理
	if result.Error != nil {
		return "", result.Error
	}

	//古いセッションを更新中にする
	session.IsUpdating = true

	//更新する
	result = dbconn.Save(&session)

	//エラー処理
	if result.Error != nil {
		return "", result.Error
	}

	return new_token, nil
}

// 更新を確定する関数
func SubmitUpdate(tokenid string, useragent string, ipaddr string) error {
	//データベース接続
	dbconn := database.GetConn()

	//新しいセッション取得
	new_session, err := GetSessionByTokenID(tokenid)

	//エラー処理
	if err != nil {
		return err
	}

	//古いセッション取得
	old_session, err := GetSessionByTokenID(new_session.UpdateID)

	//エラー処理
	if err != nil {
		return err
	}

	//更新中か
	if !old_session.IsUpdating {
		return errors.New("session is not updating")
	}

	//古いセッションを削除する
	result := dbconn.Unscoped().Delete(&old_session)

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	//新しいセッションを更新する
	new_session.UpdateID = ""
	new_session.Type = "access"
	new_session.UserAgent = useragent
	new_session.IPAddress = ipaddr
	new_session.Exp = GetExp()

	//更新する
	result = dbconn.Save(&new_session)

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// セッション取得
func GetSessionByTokenID(tokenid string) (*database.Session, error) {
	//データベース接続
	dbconn := database.GetConn()

	//セッション取得
	var session database.Session

	//セッション取得
	result := dbconn.Where(database.Session{
		TokenID: tokenid,
	}).First(&session)

	//エラー処理
	if result.Error != nil {
		return nil, result.Error
	}

	return &session, nil
}

// セッション削除
func DeleteSession(sessionid string) error {
	//データベース接続
	dbconn := database.GetConn()

	result := dbconn.Where(database.Session{
		SessionID: sessionid,
	}).Unscoped().Delete(&database.Session{})

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	return nil
}

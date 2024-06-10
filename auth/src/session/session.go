package session

import (
	"auth/model"
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
	model.Init()

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
	//古いセッション取得
	for {
		//3秒待機
		time.Sleep(time.Second * 3)

		//古いセッションを取得
		sessions, err := model.GetExperiedSession(noexp)

		//エラー処理
		if err != nil {
			log.Println(err)
			continue
		}
		

		//古いセッションを削除
		for _, session := range sessions {
			//有効期限が設定されていないとき
			if session.Exp == noexp {
				continue
			}

			//アクセストークンの場合
			if session.Type == "access" {
				//セッション取得
				get_session, err := model.GetSessionByTokenID(session.UpdateID)

				//エラー処理
				if err != nil {
					log.Println(err)
					continue
				}

				//セッション削除
				err = model.DeleteSession(get_session.SessionID)

				//エラー処理
				if err != nil {
					log.Println(err)
					continue
				}

				continue
			}

			//更新元セッション取得
			old_session, err := model.GetSessionByTokenID(session.UpdateID)

			//エラー処理
			if err == nil {
				log.Println(err)

				//更新中を外す
				old_session.IsUpdating = false

				//更新
				err := model.UpdateSession(old_session)

				//エラー処理
				if err != nil {
					log.Println(err)
					continue
				}

			} else {
				log.Println(err)
			}

			//セッション削除
			err = model.DeleteSession(session.SessionID)

			//エラー処理
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

// セッションを作成する (トークンを返す)
func GenSession(bindid string, useragent string, ipaddr string) (string, error) {
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
	session_data := model.Session{
		SessionID:  SessionID,
		UserID:     bindid,
		TokenID:    tokenid,
		UserAgent:  useragent,
		IPAddress:  ipaddr,
		IsUpdating: false,
		Exp:        GetExp(),
	}

	//セッション作成
	err = model.CreateSession(&session_data)

	//エラー処理
	if err != nil {
		return "", err
	}

	return stoken, nil
}

// 更新用のトークンを返す
// セッションを更新する
func UpdateSession(tokenid string, useragent string, ipaddr string) (string, error) {
	//セッション取得
	session, err := model.GetSessionByTokenID(tokenid)

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
	//セッションID取得
	SessionID := GenID()

	//セッション作成 (有効期限 5分)
	session_data := model.Session{
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
	err = model.CreateSession(&session_data)

	//エラー処理
	if err != nil {
		return "", err
	}

	//古いセッションを更新中にする
	session.IsUpdating = true

	//更新する
	err = model.UpdateSession(session)

	//エラー処理
	if err != nil {
		return "", err
	}

	return new_token, nil
}

// 更新を確定する関数
func SubmitUpdate(tokenid string, useragent string, ipaddr string) error {
	//新しいセッション取得
	new_session, err := model.GetSessionByTokenID(tokenid)

	//エラー処理
	if err != nil {
		return err
	}

	//古いセッション取得
	old_session, err := model.GetSessionByTokenID(new_session.UpdateID)

	//エラー処理
	if err != nil {
		return err
	}

	//更新中か
	if !old_session.IsUpdating {
		return errors.New("session is not updating")
	}

	//古いセッションを削除する
	err = model.DeleteSession(old_session.SessionID)

	//エラー処理
	if err != nil {
		return err
	}

	//新しいセッションを更新する
	new_session.UpdateID = ""
	new_session.Type = "access"
	new_session.UserAgent = useragent
	new_session.IPAddress = ipaddr
	new_session.Exp = GetExp()

	//更新する
	err = model.UpdateSession(new_session)

	//エラー処理
	if err != nil {
		return err
	}

	return nil
}
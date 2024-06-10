package session

import (
	"auth/model"
	"log"
	"time"
)

var (
	//更新用トークン有効期限 (5分)
	refresh_exp int64 = 300

	//有効期限無限
	noexp int64 = -1

	//JWT鍵
	JwtSecret string = ""
)

func Init(secret string) {
	//データベース接続
	model.Init()

	//JWT鍵の読み込み
	JwtSecret = secret
	
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

		//現在時間取得
		now := time.Now().Unix()

		//古いセッションを取得
		sessions, err := model.GetExperiedSession(now)

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

package session

import (
	"auth/database"
	"crypto/ed25519"
	"errors"
	"log"
	"time"
)

var (
	//秘密鍵
	priv_key ed25519.PrivateKey = nil

	//公開鍵
	pub_key ed25519.PublicKey = nil

	//無期限
	noexp int64 = -1
)

func Init() {
	//データベース接続
	database.Init()

	//秘密鍵と公開鍵を生成
	gen_priv,gen_pub, err := GenKey("ed25519","ed25519.pub")

	//エラー処理
	if err != nil {
		log.Println(err)
		return
	}

	//グローバル変数に格納
	priv_key = gen_priv
	pub_key = gen_pub

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
		sessions := []database.Session{}

		//db検索
		result := dbconn.Where("exp < ? AND exp <> " + string(noexp),time.Now()).Find(sessions)

		//エラー処理
		if result.Error != nil {
			log.Println(result.Error)
			return
		}

		//古いセッションを削除
		for _,session := range sessions {
			log.Println(session.Exp)

			//有効期限が設定されていないとき
			if session.Exp == noexp {
				continue
			}

			//セッション削除
			dbconn.Unscoped().Delete(&session)
		}
	}
}

//セッションを作成する (トークンを返す)
func GetSession(bindid string,useragent string,ipaddr string) (string,error) {
	//データベース接続
	dbconn := database.GetConn()

	//トークンID生成
	tokenid := GenID()

	//トークン取得
	stoken,err := GenToken(tokenid)

	//エラー処理
	if err != nil {
		return "",err
	}

	//セッションID取得
	SessionID := GenID()

	//セッション作成
	session_data := database.Session{
		SessionID: SessionID,
		UserID: bindid,
		TokenID: tokenid,
		UserAgent: useragent,
		IPAddress: ipaddr,
		IsUpdating: false,
	}

	//セッション作成
	result := dbconn.Create(&session_data)

	//エラー処理
	if result.Error != nil {
		return "",result.Error
	}

	return stoken,nil
}

//更新用のトークンを返す
//セッションを更新する
func UpdateSession(tokenid string,useragent string,ipaddr string) (string,error) {
	//データベース接続
	//dbconn := database.GetConn()

	//セッション取得
	session,err := GetSessionByTokenID(tokenid)

	//エラー処理
	if err != nil {
		return "",err
	}

	//更新用のセッションの場合
	if session.Type == "refresh" {
		return "",errors.New("session is refresh")
	}

	//更新中の場合
	if (session.IsUpdating) {
		return "",errors.New("session is updating")
	}

	//トークンID生成
	new_tokenid := GenID()
	//更新用トークン生成
	new_token,err := GenToken(new_tokenid)

	//エラー処理
	if err != nil {
		return "",err
	}

	//更新用セッション作成
	//データベース接続
	dbconn := database.GetConn()

	//セッションID取得
	SessionID := GenID()

	//セッション作成 (有効期限 5分)
	session_data := database.Session{
		SessionID: SessionID,
		UserID: session.UserID,
		TokenID: new_tokenid,
		UserAgent: useragent,
		IPAddress: ipaddr,
		IsUpdating: false,
		UpdateID: tokenid,
		Type: "refresh",
		Exp: time.Now().Add(time.Minute * 5).Unix(),
	}

	//セッション作成
	result := dbconn.Create(&session_data)

	//エラー処理
	if result.Error != nil {
		return "",result.Error
	}

	//古いセッションを更新中にする
	session.IsUpdating = true

	//更新する
	result = dbconn.Save(&session)

	//エラー処理
	if result.Error != nil {
		return "",result.Error
	}

	return new_token,nil
}

//更新を確定する関数
func SubmitUpdate(tokenid string,useragent string,ipaddr string) error {
	//データベース接続
	dbconn := database.GetConn()

	//新しいセッション取得
	new_session,err := GetSessionByTokenID(tokenid)

	//エラー処理
	if err != nil {
		return err
	}

	//古いセッション取得
	old_session,err := GetSessionByTokenID(new_session.UpdateID)

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
	new_session.Exp = noexp

	//更新する
	result = dbconn.Save(&new_session)

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	return nil
}

//セッション取得
func GetSessionByTokenID(tokenid string) (*database.Session,error) {
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
		return nil,result.Error
	}

	return &session,nil
}

//セッション削除
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
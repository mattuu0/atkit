package session

import (
	"auth/database"
	"crypto/ed25519"
	"errors"
	"log"
)

var (
	//秘密鍵
	priv_key ed25519.PrivateKey = nil

	//公開鍵
	pub_key ed25519.PublicKey = nil
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
		IsUpdate: false,
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
func UpdateSession(tokenid string) (string,error) {
	//データベース接続
	//dbconn := database.GetConn()

	//セッション取得
	session,err := GetSessionByTokenID(tokenid)

	//エラー処理
	if err != nil {
		return "",err
	}

	//更新中の場合
	if (session.IsUpdate) {
		return "",errors.New("session is updating")
	}


	GenToken(session.TokenID)

	return "",nil
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
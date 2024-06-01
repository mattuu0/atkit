package session

import "auth/database"

func Init() {
	//データベース接続
	database.Init()
}

func GetSession(bindid string,useragent string,ipaddr string) (string,error) {
	//データベース接続
	dbconn := database.GetConn()

	//セッションID取得
	SessionID := GenID()

	//セッション作成
	session_data := database.Session{
		SessionID: SessionID,
		UserID: bindid,
		TokenID: "",
		UserAgent: useragent,
		IPAddress: ipaddr,
	}

	//セッション作成
	result := dbconn.Create(&session_data)

	//エラー処理
	if result.Error != nil {
		return "",result.Error
	}

	return "",nil
}
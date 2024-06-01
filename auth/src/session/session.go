package session

import "auth/database"

func Init() {
	//データベース接続
	database.Init()
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

//セッションを更新する
func UpdateSession(tokenid string) error {
	//データベース接続
	dbconn := database.GetConn()

	//セッション取得
	session,err := GetSessionByTokenID(tokenid)

	//エラー処理
	if err != nil {
		return err
	}

	
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
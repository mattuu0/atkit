package service

import (
	"auth/model"
	"auth/session"
)

func Logout(del_session *model.Session) error {	
	//セッション削除
	err := model.DeleteSession(del_session.SessionID)

	//エラー処理
	if err != nil {
		return err
	}

	return nil
}

func StartUpdateSession(now_session *model.Session,useragent string, clientip string) (string, error) {	
	//セッション更新
	update_token,err := session.UpdateSession(now_session.TokenID,useragent,clientip)

	return update_token,err
}

func SubmitUpdateSession(now_session *model.Session,useragent string, clientip string) (error) {	
	//セッション更新
	err := session.SubmitUpdate(now_session.TokenID,useragent,clientip)

	return err
}
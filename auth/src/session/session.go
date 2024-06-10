package session

import (
	"auth/model"
	"auth/util"
	"errors"
	"time"
)

// セッションを作成する (トークンを返す)
func GenSession(bindid string, useragent string, ipaddr string) (string, error) {
	//トークンID生成
	tokenid := util.GenID()

	//トークン取得
	stoken, err := GenToken(tokenid)

	//エラー処理
	if err != nil {
		return "", err
	}

	//セッションID取得
	SessionID := util.GenID()

	//セッション作成
	session_data := model.Session{
		SessionID:  SessionID,
		UserID:     bindid,
		TokenID:    tokenid,
		UserAgent:  useragent,
		IPAddress:  ipaddr,
		IsUpdating: false,
		Exp:        util.GetExp(),
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
	get_session, err := model.GetSessionByTokenID(tokenid)

	//エラー処理
	if err != nil {
		return "", err
	}

	//更新用のセッションの場合
	if get_session.Type == "refresh" {
		return "", errors.New("session is refresh")
	}

	//更新中の場合
	if get_session.IsUpdating {
		return "", errors.New("session is updating")
	}

	//トークンID生成
	new_tokenid := util.GenID()
	//更新用トークン生成
	new_token, err := GenToken(new_tokenid)

	//エラー処理
	if err != nil {
		return "", err
	}

	//更新用セッション作成
	//セッションID取得
	SessionID := util.GenID()

	//セッション作成 (有効期限 5分)
	session_data := model.Session{
		SessionID:  SessionID,
		UserID:     get_session.UserID,
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
	get_session.IsUpdating = true

	//更新する
	err = model.UpdateSession(get_session)

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
	new_session.Exp = util.GetExp()

	//更新する
	err = model.UpdateSession(new_session)

	//エラー処理
	if err != nil {
		return err
	}

	return nil
}
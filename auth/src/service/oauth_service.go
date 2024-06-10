package service

import (
	"auth/model"
	"auth/session"
	"auth/util"
	"os"

	"github.com/markbates/goth"
	"gorm.io/gorm"
)

//トークンを返す
func Callback_Oauth(user goth.User,useragent string,clientip string) (string,error) {
	//ユーザーを取得する
	usr, err := model.GetUser(user.Provider, user.UserID)

	if err == gorm.ErrRecordNotFound {
		//見つからないときユーザーを作成する
		err := model.CreateUser(model.User{
			UserID:      user.UserID,
			Provider:    user.Provider,
			UserName:    user.Name,
			NickName:    user.NickName,
			Email:       user.Email,
			IconPath:    "./assets/UserIcons/default.png",
			ProviderUID: user.UserID,
			IsVeriry:    false,
		})

		//エラー処理
		if err != nil {
			return "",err
		}

		//ユーザーを取得
		get_usr, err := model.GetUser(user.Provider, user.UserID)

		//エラー処理
		if err != nil {
			return "",err
		}

		//アイコンパス
		//アイコンのURLが存在するか
		if user.AvatarURL != "" {
			iconpath := "./assets/UserIcons/" + get_usr.UserID + ".png"
			//アイコンを保存する
			err := util.SaveIcon(user.AvatarURL, iconpath)

			//エラー処理
			if err != nil {
				return "",err
			}

			//ユーザーを更新する
			get_usr.IconPath = iconpath

			//ユーザーを更新
			err = model.UpdateUser(get_usr)

			//エラー処理
			if err != nil {
				return "",err
			}
		}

		//ユーザーを取得
		usr = get_usr
	} else if err != nil {
		//エラー処理
		return "",err
	}


	//アイコンファイルが存在するか
	if _, err := os.Stat(usr.IconPath); err != nil {
		//存在しないとき
		//アイコンURLが存在するか
		if user.AvatarURL != "" {
			//アイコンを保存する
			err := util.SaveIcon(user.AvatarURL, usr.IconPath)

			//エラー処理
			if err != nil {
				return "",err
			}
		}
	}

	//セッションを作成する
	session_token, err := session.GenSession(usr.UserID, useragent, clientip)

	//エラー処理
	if err != nil {
		return "",err
	}

	return session_token, nil
}
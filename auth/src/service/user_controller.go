package service

import (
	"auth/model"
	"auth/util"
	"mime/multipart"
	"path/filepath"
)

//ユーザーを取得する
func GetUser(userid string) (*model.User,error) {
	return model.GetUserByID(userid)
}

func UploadIcon(userid string,icon *multipart.FileHeader) error {
	//ユーザーを取得
	user,err := model.GetUserByID(userid)

	//エラー処理
	if err != nil {
		return err
	}

	//ファイルパス指定
	savepath := filepath.Join("./assets/UserIcons",userid + ".png")

	//アイコンファイルを開く
	ofile,err := icon.Open()

	//エラー処理
	if err != nil {
		return err
	}

	//アイコンファイルを保存
	err = util.Resizeio(ofile,savepath)

	//エラー処理
	if err != nil {
		return err
	}

	//ユーザーを取得
	user.IconPath = savepath
	
	//ユーザー更新
	err = model.UpdateUser(user)

	return err
}
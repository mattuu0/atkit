package oauth

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"

	"net/http"
	"os"

	"golang.org/x/image/draw"
)


func Resizeio(iofile io.Reader,path string) error {
	//画像を開く
	file, err := os.Create(path)

	//エラー処理
	if err != nil {
		return err
	}

	defer file.Close()

	//画像をリサイズする
	//画像を読み込む
	img, _, err := image.Decode(iofile)

	//エラー処理
	if err != nil {
		return err
	}

	//リサイズ
	rimg := ResizeImage(img, 256, 256)

	//画像を保存する
	err = png.Encode(file,rimg)

	//エラー処理
	if err != nil {
		return err
	}

	return nil
}


func SaveIcon(url, path string) error {
	//リクエストを飛ばす
    response, err := http.Get(url)

	//エラー処理
    if err != nil {
        return err
    }

    defer response.Body.Close()

	//リサイズして保存
	return Resizeio(response.Body, path)
}

func ResizeImage(img image.Image, width, height int) image.Image {
	// 欲しいサイズの画像を新しく作る
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))

	// サイズを変更しながら画像をコピーする
	draw.BiLinear.Scale(newImage, newImage.Bounds(), img, img.Bounds(), draw.Over, nil)

	return newImage
}

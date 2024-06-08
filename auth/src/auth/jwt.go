package auth

import (
	"crypto/ed25519"
	"log"
	"os"
)

var (
	priv_key ed25519.PrivateKey = nil
	pub_key  ed25519.PublicKey  = nil
)

func Init() {
	//秘密鍵と公開鍵を生成
	gen_priv, gen_pub, err := GenKey(os.Getenv("Ed25519_KEY_PATH"), os.Getenv("Ed25519_PUB_PATH"))

	//エラー処理
	if err != nil {
		log.Println(err)
		return
	}

	//グローバル変数に格納
	priv_key = gen_priv
	pub_key = gen_pub
}
package authsdk

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

func ReadKey(pub_path string) (ed25519.PublicKey,error) {
	//公開鍵を読み込む
	pub_pem_bin, err := ioutil.ReadFile(pub_path)
	if err != nil {
		return nil,err
	}

	// public key デコード
	pub_block, _ := pem.Decode(pub_pem_bin)
	if pub_block == nil {
		return nil,err
	}

	//鍵をx509に変換
	parse_pub, err := x509.ParsePKIXPublicKey(pub_block.Bytes)
	if err != nil {
		return nil,err
	}

	return parse_pub.(ed25519.PublicKey),nil
}
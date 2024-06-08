package auth

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
)

// GenerateSaveEd25519 generates and saves ed25519 keys to disk after
// encoding into PEM format
func GenKey(priv_path, pub_path string) (ed25519.PrivateKey,error) {
	//鍵を読み込む
	priv_key, err := ReadKeys(priv_path)

	//成功したとき
	if err == nil {
		return priv_key,nil
	}

	//鍵を生成する
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Println(err)
		return nil,err
	}

	//Private Key をx509に変換
	priv_bin, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil,err
	}

	// private key
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: priv_bin,
	}

	//ファイルに書き込み
	err = ioutil.WriteFile(priv_path, pem.EncodeToMemory(block), 0600)
	if err != nil {
		return nil,err
	}

	// public key
	pub_bin, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil,err
	}

	// public key in PEM format
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pub_bin,
	}

	//ファイルに書き込み
	err = ioutil.WriteFile(pub_path, pem.EncodeToMemory(block), 0644)

	//エラー処理
	if err != nil {
		return nil,err
	}

	return priv,nil
}

func ReadKeys(priv_path string) (ed25519.PrivateKey,error) {
	//秘密鍵を読み込む
	priv_pem_bin, err := ioutil.ReadFile(priv_path)
	if err != nil {
		return nil,err
	}

	// private key デコード
	priv_block, _ := pem.Decode(priv_pem_bin)
	if priv_block == nil {
		return nil,err
	}

	//鍵をx509に変換
	parse_priv, err := x509.ParsePKCS8PrivateKey(priv_block.Bytes)
	if err != nil {
		return nil,err
	}

	return parse_priv.(ed25519.PrivateKey) ,nil
}
package session

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
)

// GenerateSaveEd25519 generates and saves ed25519 keys to disk after
// encoding into PEM format
func GenKey(priv_path, pub_path string) (ed25519.PrivateKey,ed25519.PublicKey,error) {
	//鍵を読み込む
	priv_key, pub_key, err := ReadKeys(priv_path, pub_path)

	//成功したとき
	if err == nil {
		return priv_key, pub_key,nil
	}

	//鍵を生成する
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Println(err)
		return nil,nil,err
	}

	//Private Key をx509に変換
	priv_bin, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil,nil,err
	}

	// private key
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: priv_bin,
	}

	//ファイルに書き込み
	err = ioutil.WriteFile(priv_path, pem.EncodeToMemory(block), 0600)
	if err != nil {
		return nil,nil,err
	}

	// public key
	pub_bin, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil,nil,err
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
		return nil,nil,err
	}

	return priv,pub,nil
}

func ReadKeys(priv_path, pub_path string) (ed25519.PrivateKey,ed25519.PublicKey,error) {
	//秘密鍵を読み込む
	priv_pem_bin, err := ioutil.ReadFile(priv_path)
	if err != nil {
		return nil,nil,err
	}

	// private key デコード
	priv_block, _ := pem.Decode(priv_pem_bin)
	if priv_block == nil {
		return nil,nil,err
	}

	//鍵をx509に変換
	parse_priv, err := x509.ParsePKCS8PrivateKey(priv_block.Bytes)
	if err != nil {
		return nil,nil,err
	}

	//公開鍵を読み込む
	pub_pem_bin, err := ioutil.ReadFile(pub_path)
	if err != nil {
		return nil,nil,err
	}

	// public key デコード
	pub_block, _ := pem.Decode(pub_pem_bin)
	if pub_block == nil {
		return nil,nil,err
	}

	//鍵をx509に変換
	parse_pub, err := x509.ParsePKIXPublicKey(pub_block.Bytes)
	if err != nil {
		return nil,nil,err
	}

	return parse_priv.(ed25519.PrivateKey),parse_pub.(ed25519.PublicKey),nil
}
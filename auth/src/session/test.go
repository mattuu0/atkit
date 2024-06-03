package session

import (
	"log"
)

func main() {
	//初期化
	Init()

	Init()
	//セッション作成
	first_token, _ := GetSession("aiueo", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.142.86 Safari/537.36", "127.0.0.1")

	log.Println(first_token)

	//トークン検証
	firstid, _ := VerifyToken(first_token)

	log.Println(firstid)

	//更新開始
	second_token, _ := UpdateSession(firstid, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.142.86 Safari/537.36", "127.0.0.1")

	log.Println(second_token)

	//トークン検証
	second_id, _ := VerifyToken(second_token)

	log.Println(second_id)

	//更新確定
	err := SubmitUpdate(second_id, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.142.86 Safari/537.36", "127.0.0.1")

	if err != nil {
		log.Println(err)
	}

	//2回目
	//更新開始
	second_token, _ = UpdateSession(second_id, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.142.86 Safari/537.36", "127.0.0.1")

	log.Println(second_token)

	//トークン検証
	second_id, _ = VerifyToken(second_token)

	log.Println(second_id)

	//更新確定
	err = SubmitUpdate(second_id, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.142.86 Safari/537.36", "127.0.0.1")

	if err != nil {
		log.Println(err)
	}

	//3回目
	//更新開始
	second_token, _ = UpdateSession(second_id, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.142.86 Safari/537.36", "127.0.0.1")

	log.Println(second_token)

	//トークン検証
	second_id, _ = VerifyToken(second_token)

	log.Println(second_id)

	//更新を開始する (失敗する)
	second_token, _ = UpdateSession(second_id, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.142.86 Safari/537.36", "127.0.0.1")

	log.Println("failed : ", second_token)

	//更新確定
	err = SubmitUpdate(second_id, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.142.86 Safari/537.36", "127.0.0.1")

	if err != nil {
		log.Println(err)
	}
	//GenerateSaveEd25519("ed25519")
}
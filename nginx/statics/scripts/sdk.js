
async function GetUser() {
    //ユーザー取得
    const req = await fetch('/auth/authed/GetUser');

    //失敗したとき
    if (req.status != 200) {
        console.log("failed to get user")
        return null;
    }

    //Jsonに変換
    const res = await req.json();
    
    return res;
}

async function Logout() {
    //ログアウト
    const req = await fetch('/auth/authed/Logout',{
        method: 'POST',
    })

    //Jsonに変換
    const res = await req.json();
    
    return res["success"];
}

//セッション更新
async function Update() {
    //更新を開始する
    const req = await fetch('/auth/authed/Update',{
        method: 'POST',
    })

    //失敗したとき
    if (req.status != 200) {
        console.log("failed to update")
        return false;
    }

    //更新を確定する
    const submit_req = await fetch('/auth/authed/SubmitUpdate',{
        method: 'POST',
    })

    //失敗したとき
    if (submit_req.status != 200) {
        console.log("failed to submit update")
        return false;
    }

    return true;
}

//アイコンを変更する関数
async function UploadIcon(ufile) {
    //送信するデータ作成
    const updata = new FormData();
    updata.append("icon",ufile);

    //送信
    const req = await fetch("/auth/uicon/upicon",{
        method: 'POST',
        body: updata
    })

    //成功したか
    if (req.status != 200) {
        return false;
    }

    return true;
}

//アイコンを取得する関数
async function GetIcon(userid) {
    return "/auth/uicon/" + userid;
}

//アクセストークン取得
async function GetToken() {
    const req = await fetch('/auth/authed/GenToken');

    //失敗したとき
    if (req.status != 200) {
        console.log("failed to get access token")
        return null;
    }

    //Jsonに変換
    const res = await req.json();
    
    return res["token"];
}
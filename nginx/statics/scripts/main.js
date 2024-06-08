const usericon = document.getElementById("usericon");

async function SetIcon(userid) {
    usericon.src = await GetIcon(userid) + "?nowtime=" + new Date().getTime();
}

const updateing_button = document.getElementById("updateing_button");

updateing_button.addEventListener('pointerdown', () => {
    const intervalId = setInterval(Update, 50)
  
    // document要素にイベント登録することで、クリックした後ボタンから動かしてもOK
    // once: true を指定して一度発火したらイベントを削除する
    document.addEventListener('pointerup', () => {        
      clearInterval(intervalId)
    }, { once: true })
})

const uicon_upload = document.getElementById("uicon_upload");

uicon_upload.addEventListener('change',async () => {
    UploadIcon(uicon_upload.files[0]);
});

var access_token = "";

async function AccessToken() {
    const token = await GetToken();

    if (token == null) {
        return;
    }

    access_token = token;
}

async function AuthTest() {
    const req = await fetch("/app/authed", {
        headers: {
            "Authorization": "Bearer " + access_token
        },
        method: "POST"
    })

    //失敗したとき
    if (req.status != 200) {
        return;
    }

    const res = await req.json();

    console.log(res);
}   

async function main() {
    const user = await GetUser();

    //失敗したとき
    if (user == null) {
        return;
    }

    console.log(user);

    SetIcon(user["UserID"]);
}

main();

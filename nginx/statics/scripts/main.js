async function main() {
    const req = await fetch('/auth/authed/GetUser')

    const res = await req.json()
    console.log(res)
}

main();

async function Logout() {
    const req = await fetch('/auth/authed/Logout',{
        method: 'POST',
    })

    const res = await req.json()
    console.log(res)
}

async function Update() {
    const req = await fetch('/auth/authed/Update',{
        method: 'POST',
    })

    const res = await req.json()
    console.log(res)

    const submit_reqq = await fetch('/auth/authed/SubmitUpdate',{
        method: 'POST',
    })

    const submit_res = await submit_reqq.json()
    console.log(submit_res)
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
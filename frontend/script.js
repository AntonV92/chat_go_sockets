let socket
let myId
let selectedUserId
let selectedUserName
let messagesList = document.getElementById("messages-list")

document.addEventListener("DOMContentLoaded", () => {
    socket = new WebSocket("ws://localhost:8000/ws");

    socket.onopen = function (e) {
        console.log('websocket connection is ok')
    };

    socket.onmessage = onMessageEvent

    socket.onclose = function (e) {
        console.log("connection closed")
    };

    socket.onerror = function (error) {
        console.log(error)
    };

});

function sendMessage() {
    let message = document.getElementById("messageText")

    if (message.value != "" && selectedUserId != undefined) {


        clientMessage = {
            from_user: myId,
            to_user: selectedUserId,
            content: message.value
        }

        date = new Date()

        socket.send(JSON.stringify(clientMessage))
        messagesList.innerHTML = prepareMessageEl(message.value, date.toLocaleString("ru", {}), 'my-message') + messagesList.innerHTML
    }
    message.value = ""
}

function selectUser(e) {
    userID = e.dataset.userid
    userName = e.dataset.username
    selectedUserId = userID

    let chatAbout = document.getElementById("chat-about")
    chatAbout.innerHTML = `<h6 class="m-b-0">${userName}</h6>`
}

function prepareMessageEl(text, date, classType) {
    return `<li class="clearfix">
    <div class="message-data">
      <span class="message-data-time">${date}</span>
    </div>
    <div class="message ${classType}">${text}</div>
  </li>`
}

function onMessageEvent(e) {
    let messageObj = JSON.parse(e.data)

    if (messageObj.type == "message") {

        let date = new Date()
        messagesList.innerHTML = prepareMessageEl(messageObj.content, date.toLocaleString("ru", {}), 'other-message') + messagesList.innerHTML

    }

    if (messageObj.type == "init") {
        myId = messageObj.content
    }

    if (messageObj.type == "users_event") {
        if (messageObj.users_online != undefined) {
            let listHtml = ""
            let usersList = document.getElementById("users-online")
            for (let [id, userObj] of Object.entries(messageObj.users_online)) {

                if (id != myId) {
                    listHtml += `<li class="clearfix online-user" data-userid="${id}" data-username="${userObj.name}" onclick="selectUser(this)">
          <img src="frontend/icon.png" alt="avatar">
          <div class="about">
            <div class="name">${userObj.name}</div>
          </div>
          </li>`
                }

            }

            usersList.innerHTML = listHtml
        }
    }

}
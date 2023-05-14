let socket
let myId
let selectedUserId
let selectedUserName

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

        console.log("to user id: " + selectedUserId)
        clientMessage = {
            from_user: myId,
            to_user: selectedUserId,
            content: message.value
        }

        socket.send(JSON.stringify(clientMessage))
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

function onMessageEvent(e) {
    let messageObj = JSON.parse(e.data)

    console.log(e.data)

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
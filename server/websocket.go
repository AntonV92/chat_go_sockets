package server

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"ws_app/user"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {

		cookie, err := r.Cookie(user.AuthCookieName)

		if err != nil {
			log.Println(err)
			return false
		}

		data := strings.Split(cookie.Value, "|")

		userID, err := strconv.Atoi(data[0])
		if err != nil {
			log.Println(err)
			return false
		}
		token := data[1]

		return user.CheckToken(userID, token)
	},
}

func getConnection() httpHanlder {
	return func(w http.ResponseWriter, r *http.Request) {

		_, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
	}
}

// func WsHandler(w http.ResponseWriter, r *http.Request) {

// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Print("upgrade:", err)
// 		return
// 	}

// 	userId, strConvErr := strconv.Atoi(r.URL.Query().Get("user_id"))

// 	defer func() {
// 		delete(ClientsOnline.ClientsList, userId)
// 		ClientsEvents <- true
// 	}()

// 	if strConvErr != nil {
// 		fmt.Println(strConvErr)
// 		return
// 	}

// 	loggedUser, isLoggedUser := user.UsersStorage[userId]
// 	if !isLoggedUser {
// 		conn.WriteMessage(websocket.TextMessage, []byte("Login session is expired"))
// 		return
// 	}

// 	ClientsOnline.ClientsList[loggedUser.Id] = loggedUser.Name
// 	loggedUser.WsConn = *conn
// 	ClientsEvents <- true

// 	fmt.Printf("Connected: %s\n", loggedUser.Name)

// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("read ws message error: ", err)
// 			break
// 		}

// 		sendMessageErr := sendClientMessage(message)
// 		if sendMessageErr != nil {
// 			fmt.Printf("send client message error: %v\n", sendMessageErr)
// 			continue
// 		}
// 	}
// }

// func sendClientMessage(message []byte) error {

// 	mes := ClientMessage{}
// 	decodeError := json.Unmarshal(message, &mes)
// 	if decodeError != nil {
// 		return decodeError
// 	}
// 	sendError := user.UsersStorage[mes.ToUserId].WsConn.WriteMessage(websocket.TextMessage, []byte(message))
// 	if sendError != nil {
// 		return sendError
// 	}

// 	return nil
// }

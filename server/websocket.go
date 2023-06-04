package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"ws_app/helpers"
	"ws_app/user"

	"github.com/gorilla/websocket"
)

const (
	MessageTypeUsersEvent    = "users_event"
	MessageTypeSimpleMessage = "message"
	MessageTypeInitMessage   = "init"
)

type Message struct {
	Type        string             `json:"type"`
	Content     string             `json:"content,omitempty"`
	UsersOnline map[int]*user.User `json:"users_online,omitempty"`
}

type ClientMessage struct {
	Type     string `json:"type"`
	FromUser int    `json:"from_user,string"`
	ToUser   int    `json:"to_user,string"`
	Content  string `json:"content"`
}

func getConnection() httpHanlder {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie(user.AuthCookieName)
		userID, token := ParseAuthCookie(cookie.Value)

		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {

				if err != nil {
					log.Println(err)
					return false
				}
				return user.CheckToken(userID, token)
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		userInstance, err := user.GetUserById(userID)
		if err != nil {
			log.Println(err)
		}
		userInstance.WsConn = conn
		user.UsersOnline[userID] = userInstance
		user.ConnectionEvents <- true

		defer func() {
			delete(user.UsersOnline, userID)
			user.ConnectionEvents <- true
			fmt.Printf("User: %d disconnected\n", userID)
		}()

		initMes := Message{
			Type:    MessageTypeInitMessage,
			Content: strconv.Itoa(userID),
		}
		jsonInit, err := json.Marshal(initMes)
		if err != nil {
			log.Println(err)
			return
		}

		// send init data to client
		conn.WriteMessage(websocket.TextMessage, jsonInit)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("Get message from client: %s\n", message)
			processClientMessage(message)
		}
	}
}

func processClientMessage(message []byte) {
	var clientMes = ClientMessage{
		Type: MessageTypeSimpleMessage,
	}
	err := json.Unmarshal(message, &clientMes)
	if err != nil {
		log.Println(err)
		return
	}

	clientMes.Content = helpers.PrepareTextLinks(clientMes.Content)

	mes, err := json.Marshal(clientMes)
	if err != nil {
		log.Println(err)
		return
	}

	targetClient := user.UsersOnline[clientMes.ToUser]

	if targetClient == nil {
		return
	}

	sendErr := targetClient.WsConn.WriteMessage(websocket.TextMessage, mes)
	if sendErr != nil {
		log.Println(sendErr)
		return
	}
}

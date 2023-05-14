package server

import (
	"fmt"
	"log"
	"net/http"
	"ws_app/user"

	"github.com/gorilla/websocket"
)

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

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("Get message from client: %s\n", message)
		}
	}
}

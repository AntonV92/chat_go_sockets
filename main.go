package main

import (
	"encoding/json"
	"fmt"
	"log"
	"ws_app/db"
	"ws_app/server"
	"ws_app/user"

	"github.com/gorilla/websocket"
)

func main() {
	go broadcaster()
	db.InitDbConnection()
	server.InitServer()

	defer db.DbConn.Close()
}

func broadcaster() {
	for {
		select {
		case <-user.ConnectionEvents:

			for _, onlineUser := range user.UsersOnline {
				fmt.Println(onlineUser.Name)

				message := server.Message{
					Type:        server.MessageTypeUsersEvent,
					UsersOnline: user.UsersOnline,
				}

				jsonMessage, err := json.Marshal(message)
				if err != nil {
					log.Println(err)
					continue
				}

				onlineUser.WsConn.WriteMessage(websocket.TextMessage, jsonMessage)
			}
		}
	}
}

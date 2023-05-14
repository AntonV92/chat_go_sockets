package main

import (
	"fmt"
	"ws_app/db"
	"ws_app/server"
	"ws_app/user"
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
			}
		}
	}
}

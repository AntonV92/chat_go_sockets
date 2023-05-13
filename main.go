package main

import (
	"ws_app/db"
	"ws_app/server"
)

func main() {
	db.InitDbConnection()
	server.InitServer()

	defer db.DbConn.Close()
}

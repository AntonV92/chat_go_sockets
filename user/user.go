package user

import (
	"database/sql"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	Id           int            `json:"id"`
	Name         string         `json:"name"`
	Password     string         `json:"password"`
	Token        sql.NullString `json:"token"`
	Token_update time.Time      `json:"token_update"`
	WsConn       websocket.Conn
}

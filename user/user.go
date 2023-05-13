package user

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
	"ws_app/db"
	"ws_app/helpers"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

const (
	chars               = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	TokenExpiredMinutes = 60
	AuthCookieName      = "authCookie"
)

type User struct {
	Id           int            `json:"id"`
	Name         string         `json:"name"`
	Password     string         `json:"password"`
	Token        sql.NullString `json:"token"`
	Token_update time.Time      `json:"token_update"`
	WsConn       websocket.Conn
}

type CookieValue struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

// get hashed string from given password
func PasswordHash(pass string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 0)
	if err != nil {
		log.Fatal(err)
	}

	return string(hash)
}

func GenerateToken(lenght int) string {

	token := ""

	for i := 0; i < lenght; i++ {
		randIndex := rand.Intn(len(chars) - 1)
		token += string(chars[randIndex])
	}

	return token
}

func CheckToken(userId int, token string) bool {
	db := db.DbConn

	var currentToken string
	row := db.QueryRow("SELECT token FROM users WHERE id = $1 LIMIT 1;", userId)
	row.Scan(&currentToken)

	return token == currentToken
}

func Login(login string, pass string) (User, error) {
	db := db.DbConn

	user := User{}

	userRecord := db.QueryRow("SELECT * FROM users WHERE name = $1 LIMIT 1;", login)
	userRecord.Scan(&user.Id, &user.Name, &user.Password, &user.Token, &user.Token_update)

	if !CheckPassword(pass, user.Password) {
		return User{}, fmt.Errorf("Wrong password")
	}

	tokenCreated := helpers.GetTimeDiffNow(user.Token_update).Minutes()

	var token sql.NullString
	if user.Token.String == "" || int(tokenCreated) > TokenExpiredMinutes {
		genToken := GenerateToken(30)
		updatedDate := time.Now().Format(time.DateTime)
		db.QueryRow("UPDATE users SET token = $1, token_update = $2 WHERE id = $3 returning token;",
			genToken, updatedDate, user.Id).Scan(&token)
		user.Token_update = time.Now()
	} else {
		db.QueryRow("SELECT token FROM users WHERE id = $1 LIMIT 1;", user.Id).Scan(&token)
	}

	user.Token = token
	return user, nil
}

// compare password with string hash
func CheckPassword(pass string, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(pass))
	if err != nil {
		return false
	}

	return true
}

package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DbConn *sql.DB

func InitDbConnection() {

	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal(envErr)
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dataSource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, pass, dbName)

	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		log.Fatal(err)
	}

	if dbConnError := db.Ping(); dbConnError != nil {
		log.Fatal(dbConnError)
	}

	DbConn = db
}

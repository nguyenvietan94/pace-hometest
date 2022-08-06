package middleware

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var dbInstance *sql.DB

// create connection with postgres db
func DbConnect() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	if dbInstance == nil {
		db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
		if err != nil {
			panic(err)
		}

		err = db.Ping()
		if err != nil {
			panic(err)
		}
		dbInstance = db
		if dbInstance == nil {
			fmt.Println("dbInstance == nil")
		}
		fmt.Println("Successfully connected!")
	}
	return dbInstance
}

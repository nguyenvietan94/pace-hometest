package middleware

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var dbInstance *sql.DB

func DbConnect() *sql.DB {
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
		log.Println("Successfully connected to database!")
	}

	return dbInstance
}

package main

import (
	"log"
	"net/http"
	"pace-hometest/middleware"
	"pace-hometest/router"
)

func main() {
	db := middleware.DbConnect()
	defer db.Close()
	r := router.Router()
	log.Println("Started server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

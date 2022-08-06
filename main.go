package main

import (
	"fmt"
	"log"
	"net/http"
	"pace-hometest/middleware"
	"pace-hometest/router"
)

func main() {

	_ = middleware.DbConnect()

	r := router.Router()
	fmt.Println("Started server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

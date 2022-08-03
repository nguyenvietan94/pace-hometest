package main

import (
	"fmt"
	"log"
	"net/http"
	"pace-hometest/router"
)

func main() {
	r := router.Router()
	fmt.Println("Started server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

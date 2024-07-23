package main

import (
	"log"
	"passkey-auth/router"
)

func main() {
	router, err := router.New()
	if err != nil {
		panic(err)
	}

	log.Println("Server started on port 8082")
	log.Fatal(router.Start(":8082"))
}

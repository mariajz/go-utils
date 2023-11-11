package main

import (
	"log"
	"os"

	service "github.com/mariajz/go-utils/mongodb-client/service"
)

func main() {
	log.Println("Initializing mongodb client")
	log.Println("MONGO_URL", os.Getenv("MONGO_URL"))

	service.InitDB()
}

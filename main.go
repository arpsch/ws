package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	route "ws/api/http"
)

func init() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func readServerPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}

func main() {
	port := readServerPort()
	router := route.SetupRoutes()
	log.Printf("Server is listening on %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}

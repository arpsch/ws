package main

import (
	"log"
	"net/http"
	"os"

	route "ws/api/http"
)

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

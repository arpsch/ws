package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"

	api_http "ws/api/http"
	"ws/app"
	mongo "ws/store/mongo"
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

func setupServer() (*httprouter.Router, error) {
	// set up mongo
	db, err := mongo.NewMongo()
	if err != nil {
		return nil, err
	}

	// set up Image collector
	icApp := app.NewImageCollector(db)

	// set up api handlers for image collector
	appHandler := api_http.NewAppHandlers(icApp)
	routes := appHandler.SetupRoutes()

	return routes, nil
}

func main() {
	router, err := setupServer()
	if err != nil {
		log.Printf(" failed to set up routes, exiting")
		return
	}

	port := readServerPort()
	log.Printf("Server is listening on %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}

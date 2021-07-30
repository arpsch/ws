package http

import (
	"net/http"
	"ws/logger"

	"github.com/julienschmidt/httprouter"
)

func SetupRoutes() *httprouter.Router {
	router := httprouter.New()
	router.Handler("GET", "/", http.FileServer(http.Dir("public/")))
	router.HandlerFunc("GET", "/hello/:name", logger.Log(Hello))
	router.HandlerFunc("POST", "/upload", logger.Log(fileUploadRequestHandler))

	return router
}

package http

import (
	"ws/logger"

	"github.com/julienschmidt/httprouter"
)

func SetupRoutes() *httprouter.Router {
	router := httprouter.New()
	//router.Handler("GET", "/", http.FileServer(http.Dir("public/")))
	router.HandlerFunc("GET", "/", logger.Log(IndexHandler))
	router.HandlerFunc("POST", "/upload", logger.Log(FileUploadRequestHandler))

	return router
}

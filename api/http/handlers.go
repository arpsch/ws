package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"ws/app"
	"ws/logger"

	"github.com/julienschmidt/httprouter"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 8 // 8MB

type appHandlers struct {
	app app.ImageCollectorApp
}

func NewAppHandlers(ic app.ImageCollectorApp) *appHandlers {
	return &appHandlers{
		app: ic,
	}
}

func (ah *appHandlers) SetupRoutes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc("GET", "/", logger.Log(ah.IndexHandler))
	router.HandlerFunc("POST", "/upload", logger.Log(ah.FileUploadRequestHandler))

	return router
}

func readAuthToken() string {
	at := os.Getenv("AUTH_TOKEN")
	if at != "" {
		return at
	}
	return "SecretToken"
}

func (ah *appHandlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set the auth cookie
	expire := time.Now().Add(3 * time.Second)
	cookie := http.Cookie{
		Name:    "auth",
		Value:   readAuthToken(),
		Expires: expire,
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	if r.URL.Path == "/" {
		http.ServeFile(w, r, "public/index.html")
	}
}

func (ah *appHandlers) FileUploadRequestHandler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// FormFile use 32MB by default
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	auth := r.FormValue("auth")

	fmt.Printf("############### AUTH TOKEN %s ##################\n", auth)

	if auth != readAuthToken() {
		http.Error(w, errors.New("authentication is failed").Error(), http.StatusForbidden)
		return
	}
	// read the file header
	_, fh, err := r.FormFile("data")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if fh == nil {
		http.Error(w, errors.New("failed to read file header").Error(), http.StatusBadRequest)
		return
	}

	// validate against set threshold for the file
	if fh.Size > MAX_UPLOAD_SIZE {
		http.Error(w, fmt.Sprintf("file %s seems to larger than %d", fh.Filename, MAX_UPLOAD_SIZE), http.StatusBadRequest)
		return
	}

	err = ah.app.AddImageInformation(ctx, fh)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Upload successful")
}

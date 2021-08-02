package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 8 // 8MB

func readAuthToken() string {
	at := os.Getenv("AUTH_TOKEN")
	if at != "" {
		return at
	}
	return "SecretToken"
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
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

func FileUploadRequestHandler(w http.ResponseWriter, r *http.Request) {
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

	if auth != readAuthToken() {
		http.Error(w, errors.New("authentication is failed").Error(), http.StatusForbidden)
		return
	}
	// read the file header
	_, fh, err := r.FormFile("data")

	if fh == nil {
		http.Error(w, errors.New("failed to read file header").Error(), http.StatusBadRequest)
		return
	}

	// validate against set threshold for the file
	if fh.Size > MAX_UPLOAD_SIZE {
		http.Error(w, fmt.Sprintf("file %s seems to larger than %d", fh.Filename, MAX_UPLOAD_SIZE), http.StatusBadRequest)
		return
	}

	// read the content of uploaded file
	upFile, err := fh.Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer upFile.Close()

	buff := make([]byte, 512)
	_, err = upFile.Read(buff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// identify the content type for -JPEG or PNG
	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" {
		http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
		return
	}

	// reset the pointer back to start of the file
	_, err = upFile.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// store all the uploaded files under files folder
	err = os.MkdirAll("./files", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create a destination file - .fileNameTime
	saveFile, err := os.Create(fmt.Sprintf("./files/%s%d", filepath.Ext(fh.Filename), time.Now().UnixNano()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer saveFile.Close()

	_, err = io.Copy(saveFile, upFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Upload successful")
}

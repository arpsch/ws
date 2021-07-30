package logger

import (
	"log"
	"net/http"
	"time"
)

func Log(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v - %s\n", time.Now(), r.URL.Path)
		f(w, r)
	}
}

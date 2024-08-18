package handlers

import (
	"log"
	"net/http"
	"time"
)

func TimerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("[INFO] %s %s time: %fs", r.Method, r.URL, time.Since(start).Seconds())
		}()
		next.ServeHTTP(rw, r)
	})
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var Port int
var Dir string

func init() {
	flag.IntVar(&Port, "port", 8000, "port to listen on")
	flag.StringVar(&Dir, "dir", ".", "directory to serve")
}

// ResponseWriter wraps http.ResponseWriter to capture the HTTP status code
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func CacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

// Handler wraps http.Handler to log served files
type Handler struct {
	http.Handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := &ResponseWriter{w, 0}
	h.Handler.ServeHTTP(rw, r)
	log.Println(r.RemoteAddr, r.Method, rw.StatusCode, r.URL)
}

func main() {
	flag.Parse()
	fileServer := http.FileServer(http.Dir(Dir))
	cacheControlledHandler := CacheControlMiddleware(fileServer)
	handler := &Handler{cacheControlledHandler}
	http.Handle("/", handler)
	addr := fmt.Sprintf(":%d", Port)
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

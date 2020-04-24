package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var Port int
var Dir string
var UseCors bool

func init() {
	flag.IntVar(&Port, "port", 8000, "port to listen on")
	flag.StringVar(&Dir, "dir", ".", "directory to serve")
	flag.BoolVar(&UseCors, "cors", false, "Set true to enable 'Access-Control-Allow-Origin: *' on all requests")
}

// ResponseWriter wraps http.ResponseWriter to capture the HTTP status code
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode

	if UseCors {
		w.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
	}

	w.ResponseWriter.WriteHeader(statusCode)
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
	handler := &Handler{http.FileServer(http.Dir(Dir))}
	http.Handle("/", handler)
	addr := fmt.Sprintf(":%d", Port)
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

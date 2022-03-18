package main

import (
	"log"
	"net/http"
)

func WithLogging(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI, "----", r.URL.String())
		h.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(logFn)
}

func Display(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	log.Println(url)
	w.Write([]byte(url))
}

func main() {
	http.HandleFunc("/", Display)
	log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))
}

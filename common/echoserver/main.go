package main

import (
	"log"
	"net/http"
	"strings"
	"text/template"
)

func WithLogging(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(logFn)
}

type Resp struct {
	URL      string
	ClientIP string
	Proxy    []string
	Headers  http.Header
	Cookies  []*http.Cookie
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	var resp Resp
	resp.URL = r.URL.String()
	if r.Header.Get("X-Forwared-For") == "" {
		resp.ClientIP = strings.Split(r.RemoteAddr, ":")[0]
	} else {
		resp.ClientIP = strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
		resp.Proxy = strings.Split(r.Header.Get("X-Forwarded-For"), ",")[1:]
	}
	resp.Headers = r.Header
	resp.Cookies = r.Cookies()
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	t.Execute(w, resp)
}

func main() {
	http.HandleFunc("/", EchoHandler)
	log.Fatal(http.ListenAndServe("0..0.0:80", nil))
}

package main

import (
	"embed"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

//go:embed templates/*
var content embed.FS

func WithLogging(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(logFn)
}

type Resp struct {
	URL         string
	ClientIP    string
	Proxy       []string
	Headers     http.Header
	Cookies     []*http.Cookie
	RequestData string
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	var resp Resp
	resp.URL = r.URL.String()
	resp.ClientIP = strings.Split(r.RemoteAddr, ":")[0]
	if r.Header.Get("X-Forwared-For") != "" {
		resp.Proxy = strings.Split(r.Header.Get("X-Forwarded-For"), ",")
	}
	resp.Headers = r.Header
	resp.Cookies = r.Cookies()
	if r.Body != nil {
		defer r.Body.Close()
		rData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			resp.RequestData = err.Error()
		}
		resp.RequestData = string(rData)
	}
	t, err := template.ParseFS(content, "templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	t.Execute(w, resp)
}

func printRequestBody(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}
	log.Println(string(data))
	w.WriteHeader(http.StatusOK)
	return
}

func main() {
	http.HandleFunc("/", EchoHandler)
	http.HandleFunc("/os", systemCheck)
	log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))
}

package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func systemCheck(w http.ResponseWriter, r *http.Request) {
	environmentVariables := os.Environ()
	data, err := json.Marshal(environmentVariables)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Write(data)
}

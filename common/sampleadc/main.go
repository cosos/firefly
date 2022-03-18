package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	admission "k8s.io/api/admission/v1"
	k8smeta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AdmissionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		log.Printf("No request body has been recieved")
		http.Error(w, fmt.Sprintf("No request body has been recieved"), http.StatusBadRequest)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusBadRequest)
	}
	if len(body) == 0 {
		log.Printf("request body is empty")
		http.Error(w, fmt.Sprintf("request body is empty"), http.StatusBadRequest)
	}
	log.Println(string(body))
	reviewRequest := admission.AdmissionReview{}
	if err := json.Unmarshal(body, &reviewRequest); err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, fmt.Sprintf("Error parsing request body: %v", err), http.StatusBadRequest)
	}
	result, err := CheckRequest(reviewRequest.Request)
	response := admission.AdmissionResponse{
		UID:     reviewRequest.Request.UID,
		Allowed: result,
	}
	if err != nil {
		response.Result = &k8smeta.Status{
			Message: fmt.Sprintf("%v", err),
			Reason:  k8smeta.StatusReasonForbidden,
		}
	}
	outReview := admission.AdmissionReview{
		TypeMeta: reviewRequest.TypeMeta,
		Request:  reviewRequest.Request,
		Response: &response,
	}
	json, err := json.Marshal(outReview)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response %v", err), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(json); err != nil {
			log.Printf("Error writing response %v", err)
			http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
		}
	}
}

func CheckRequest(request *admission.AdmissionRequest) (bool, error) {
	if request.Namespace == "" {
		log.Printf("no namespace has been provided for %s", request.Kind.Kind)
		return false, errors.New("no namespace has been provided")
	}
	if request.Namespace == "kube-system" {
		log.Printf("[Warning] This object will be put into kube-system namespace")
		return true, nil
	}
	return false, errors.New("unauthorized request from admission controller")
}

func main() {
	http.HandleFunc("/admission", AdmissionHandler)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong\n"))
	})
	server := &http.Server{
		Addr: "0.0.0.0:80",
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

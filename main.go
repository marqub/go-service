package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type jsonStatus struct {
	Status       string       `json:"status"`
	Name         string       `json:"name"`
	Dependencies []jsonStatus `json:"dependencies,omitempty"`
}

var servicename string
var dependency string

func main() {
	servicename = os.Getenv("SVC_NAME")
	dependency = os.Getenv("DEPENDENCY_NAME")
	router := mux.NewRouter()

	router.Methods("GET").Path("/status").Name("Status").HandlerFunc(status)

	log.Infof("Server started. Name: %s, Dependency: %s", servicename, dependency)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	globalStatus := jsonStatus{Status: "OK", Name: servicename}

	if dependency != "" {
		var client http.Client
		resp, err := client.Get(dependency + "/status")
		if err != nil {
			log.Errorf("dependency unreachable: %v ", err)
			globalStatus.Dependencies = append(globalStatus.Dependencies, jsonStatus{Status: "UNKNOWN", Name: dependency})
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var dependencyResponse jsonStatus
				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Errorf("response can't be read: %v ", err)
					returnError(http.StatusInternalServerError, "Invalid response", w)
					return
				}
				// Unmarshall the data
				if err := json.Unmarshal(data, &dependencyResponse); err != nil {
					// unprocessable entity
					log.Errorf("unprocessable entity: %v", err)
					returnError(http.StatusPreconditionFailed, "Invalid status response", w)
					return
				}
				globalStatus.Dependencies = append(globalStatus.Dependencies, dependencyResponse)
			}
		}
	}

	if err := json.NewEncoder(w).Encode(globalStatus); err != nil {
		returnError(http.StatusInternalServerError, "Invalid response", w)
	}
}

func returnError(code int, message string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(422)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: code, Text: message}); err != nil {
		log.Error(err)
	}
}

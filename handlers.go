package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This would be a general status of the system.")
}

func AddWebsite(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var website Website
	err := decoder.Decode(&website)
	if err != nil {
		// panic()
	}
	err = AddToDatabase(website)
	if err != nil {
	}
	MonitorWebsite(website)

	w.WriteHeader(204)
}

func RemoveWebsite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	websiteId := vars["id"]
	fmt.Fprintf(w, "Will remove website %s from uptime monitoring.", websiteId)
}

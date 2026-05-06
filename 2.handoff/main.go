package main

import (
	"log"
	"net/http"
)

func main() {
	memoryStore := MemoryStore{incidents: make(map[string]Incident)}
	incHandler := IncidentHandler{Store: &memoryStore}
	router := getRouter(incHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

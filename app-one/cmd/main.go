package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var Version = "dev"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/version", handleVersion)

	log.Printf("Starting app-one server version %s on port %s", Version, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	log.Printf("Stopped app-one server on port %s", port)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to App One!\n")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy"}`)
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"version":"%s"}`, Version)
}

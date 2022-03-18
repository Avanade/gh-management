package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

import (
	"log"
)

var host = flag.String("host", ":8080", "The host of the application.")

func loadEnvironmentalVariables() {
	err := godotenv.Load()
	log.Printf("Loaded as %v\n", os.Getenv("ENV"))
	if err != nil && os.Getenv("ENV") == "local" {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnvironmentalVariables()
	http.Handle("/", &templateHandler{filename: "index.html"})
	http.HandleFunc("/ossLogic", handleGraphRequest)
	port := os.Getenv("PORT")
	log.Printf("Starting server on port %v\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	target := os.Getenv("CUSTOM_DOMAIN")
	if target == "" {
		log.Fatal("CUSTOM_DOMAIN environment variable is not set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	})

	log.Println("Redirecting to", target)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

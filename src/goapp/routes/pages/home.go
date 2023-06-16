package routes

import (
	"log"
	"net/http"

	"main/pkg/template"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("APP RUNNING")
	template.UseTemplate(&w, r, "home", nil)
}

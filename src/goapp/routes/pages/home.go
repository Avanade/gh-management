package routes

import (
	"log"
	template "main/pkg/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("APP RUNNING")
	template.UseTemplate(&w, r, "home", nil)
}

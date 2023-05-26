package routes

import (
	"log"
	template "main/pkg/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("TEST : %s", "DISPLAY LOGS")
	template.UseTemplate(&w, r, "home", nil)
}

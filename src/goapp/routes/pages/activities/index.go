package routes

import (
	"log"
	"net/http"

	"main/pkg/session"
	"main/pkg/template"
)

func ActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, err := session.Store.Get(r, "auth-session")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	template.UseTemplate(&w, r, "activities/index", sessionaz)
}

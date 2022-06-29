package routes

import (
	session "main/pkg/session"
	template "main/pkg/template"
	"net/http"
)

func ActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	template.UseTemplate(&w, r, "activities/index", sessionaz)
}

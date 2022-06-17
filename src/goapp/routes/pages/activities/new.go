package routes

import (
	template "main/pkg/template"
	"net/http"
)

func ActivitiesNewHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "activities/new", nil)
}

package routes

import (
	template "main/pkg/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "home", nil)
}

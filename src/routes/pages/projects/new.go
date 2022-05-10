package routes

import (
	template "main/pkg/template"
	"net/http"
)

func ProjectsNewHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "projects/new", nil)
}

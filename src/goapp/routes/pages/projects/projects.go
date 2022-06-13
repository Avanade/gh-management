package routes

import (
	template "main/pkg/template"
	"net/http"
)

func Projects(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "projects/projects", nil)
}

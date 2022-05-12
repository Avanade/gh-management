package routes

import (
	template "main/pkg/template"
	"net/http"
)

func ProjectsNewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		template.UseTemplate(&w, r, "projects/new", nil)
	case "POST":
		w.WriteHeader(http.StatusBadRequest)
	}
}

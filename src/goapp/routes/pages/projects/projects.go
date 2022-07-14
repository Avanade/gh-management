package routes

import (
	template "main/pkg/template"
	"net/http"
)

func Projects(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "projects/projects", nil)
}

func MakePublic(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "projects/makepublic", nil)
}

package routes

import (
	session "main/pkg/session"
	template "main/pkg/template"
	"net/http"
)

func Projects(w http.ResponseWriter, r *http.Request) {
	isAdmin, _ := session.IsUserAdmin(w, r)

	data := map[string]interface{}{
		"isAdmin": isAdmin,
	}
	template.UseTemplate(&w, r, "projects/projects", data)
}

func MakePublic(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "projects/makepublic", nil)
}

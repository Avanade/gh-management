package routes

import (
	session "main/pkg/session"
	template "main/pkg/template"
	"net/http"
)

func Projects(w http.ResponseWriter, r *http.Request) {
	pageData := make(map[string]interface{})

	isAdmin, err := session.IsUserAdmin(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData["isUserAdmin"] = isAdmin

	template.UseTemplate(&w, r, "projects/projects", pageData)
}

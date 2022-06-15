package routes

import (
	session "main/pkg/session"
	template "main/pkg/template"
	"net/http"
)

func Projects(w http.ResponseWriter, r *http.Request) {
	pageData := make(map[string]interface{})

	// Check session
	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData["isUserAdmin"] = session.Values["isUserAdmin"].(bool)

	template.UseTemplate(&w, r, "projects/projects", pageData)
}

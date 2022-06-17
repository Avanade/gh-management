package routes

import (
	"main/pkg/session"
	"main/pkg/template"
	"net/http"
)

func ListCommunityMembers(w http.ResponseWriter, r *http.Request) {
	// Check if user is an admin
	isAdmin, err := session.IsUserAdmin(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isAdmin {
		template.UseTemplate(&w, r, "admin/communitymembers", nil)
	} else {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
}

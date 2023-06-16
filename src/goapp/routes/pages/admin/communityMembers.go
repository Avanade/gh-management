package routes

import (
	"net/http"

	"main/pkg/template"
)

func ListCommunityMembers(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/communitymembers", nil)
}

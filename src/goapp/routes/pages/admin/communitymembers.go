package routes

import (
	"main/pkg/template"
	"net/http"
)

func ListCommunityMembers(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/communitymembers", nil)
}

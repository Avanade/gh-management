package routes

import (
	"net/http"

	"main/pkg/template"
)

func CommunityMembersHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/communitymembers", nil)
}

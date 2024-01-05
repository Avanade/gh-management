package routes

import (
	"main/pkg/template"
	"net/http"
)

func CommunityApproversHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "/community/communityapprovers", nil)
}

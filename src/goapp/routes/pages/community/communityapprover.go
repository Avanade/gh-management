package routes

import (
	template "main/pkg/template"
	"net/http"
)

func CommunityApproverHandler(w http.ResponseWriter, r *http.Request) {

	template.UseTemplate(&w, r, "/community/communityapprovers", nil)
}

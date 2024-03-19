package routes

import (
	"main/pkg/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "otherRequests/index", nil)
}

func RequestOrganization(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "otherRequests/organization", nil)
}

func RequestGitHubCopilot(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "otherRequests/githubCopilot", nil)
}

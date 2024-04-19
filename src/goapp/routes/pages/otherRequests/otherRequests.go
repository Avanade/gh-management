package routes

import (
	"main/pkg/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "otherRequests/index", nil)
}

func RequestNewOrganization(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "otherRequests/organization", nil)
}

func RequestGitHubCopilot(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "otherRequests/githubCopilot", nil)
}

func RequestOrganizationAccess(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "otherRequests/organizationaccess", nil)
}

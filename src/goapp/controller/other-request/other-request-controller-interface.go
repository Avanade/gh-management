package otherRequest

import (
	"net/http"
)

type OtherRequestController interface {
	IndexHandler(w http.ResponseWriter, r *http.Request)
	RequestGitHubCopilot(w http.ResponseWriter, r *http.Request)
	RequestGitHubOrganization(w http.ResponseWriter, r *http.Request)
	RequestGitHubOrganizationAccess(w http.ResponseWriter, r *http.Request)
	RequestAdoOrganization(w http.ResponseWriter, r *http.Request)
}

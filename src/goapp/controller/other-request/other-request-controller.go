package otherRequest

import (
	"main/pkg/template"
	"net/http"
)

type otherRequestController struct{}

func NewOtherRequestController() OtherRequestController {
	return &otherRequestController{}
}

func (c *otherRequestController) IndexHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "otherRequests/index", nil)
}

func (c *otherRequestController) RequestGitHubCopilot(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "otherRequests/githubCopilot", nil)
}

func (c *otherRequestController) RequestGitHubOrganization(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "otherRequests/organization", nil)
}

func (c *otherRequestController) RequestGitHubOrganizationAccess(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "otherRequests/organizationaccess", nil)
}

func (c *otherRequestController) RequestAdoOrganization(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "otherRequests/adoOrganization", nil)
}

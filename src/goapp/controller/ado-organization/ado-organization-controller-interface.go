package adoOrganization

import (
	"net/http"
)

type AdoOrganizationController interface {
	CreateAdoOrganizationRequest(w http.ResponseWriter, r *http.Request)
	GetAdoOrganizationByUser(w http.ResponseWriter, r *http.Request)
	GetAdoOrganizationApprovalRequests(w http.ResponseWriter, r *http.Request)
}

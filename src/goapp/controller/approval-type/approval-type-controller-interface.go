package approvaltype

import "net/http"

type ApprovalTypeController interface {
	CreateApprovalType(w http.ResponseWriter, r *http.Request)
	GetApprovalTypeById(w http.ResponseWriter, r *http.Request)
	GetApprovalTypes(w http.ResponseWriter, r *http.Request)
}

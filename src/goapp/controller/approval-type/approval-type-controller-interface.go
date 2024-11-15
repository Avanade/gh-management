package approvaltype

import "net/http"

type ApprovalTypeController interface {
	GetApprovalTypeById(w http.ResponseWriter, r *http.Request)
	GetApprovalTypes(w http.ResponseWriter, r *http.Request)
}

package approvaltype

import "net/http"

type ApprovalTypeController interface {
	GetApprovalTypes(w http.ResponseWriter, r *http.Request)
}

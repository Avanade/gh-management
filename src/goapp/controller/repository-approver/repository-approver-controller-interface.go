package repositoryApprover

import (
	"net/http"
)

type RepositoryApproverController interface {
	GetLegalApprovers(w http.ResponseWriter, r *http.Request)
}

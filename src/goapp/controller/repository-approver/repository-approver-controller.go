package repositoryApprover

import (
	"encoding/json"
	"main/config"
	"main/service"
	"net/http"
)

type repositoryApproverController struct {
	*service.Service
	Conf config.ConfigManager
}

func NewRepositoryApproverController(service *service.Service, conf config.ConfigManager) RepositoryApproverController {
	return &repositoryApproverController{service, conf}
}

func (c *repositoryApproverController) GetLegalApprovers(w http.ResponseWriter, r *http.Request) {
	result, err := c.RepositoryApprover.Get(c.Conf.GetLegalApprovalTypeId())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

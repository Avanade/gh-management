package routes

import (
	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"
	"net/http"

	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
)

func FillOutApprovers(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	result, err := db.SelectApprovalTypes()
	if err != nil {
		logger.LogTrace(err.Error(), contracts.Error)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, approver := range result {
		if approver["ApproverUserPrincipalName"] == nil {
			continue
		}

		err = db.InsertApprover(db.Approver{
			ApprovalTypeId: int(approver["Id"].(int64)),
			ApproverEmail:  approver["ApproverUserPrincipalName"].(string),
		})
		if err != nil {
			logger.LogTrace(err.Error(), contracts.Error)
		}
	}
}

func FillOutApprovalRequestApprovers(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	projectApprovals, err := db.GetProjectApprovals()
	if err != nil {
		logger.LogTrace(err.Error(), contracts.Error)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, projectApproval := range projectApprovals {
		if projectApproval["ApproverUserPrincipalName"] == nil {
			continue
		}
		id := int(projectApproval["Id"].(int64))
		err = db.InsertApprovalRequestApprover(db.ApprovalRequestApprover{
			ApprovalRequestId: id,
			ApproverEmail:     projectApproval["ApproverUserPrincipalName"].(string),
		})
		if err != nil {
			logger.LogTrace(err.Error(), contracts.Error)
		}

		if projectApproval["ApprovalDate"] != nil && projectApproval["ApproverUserPrincipalName"] != nil {
			err = db.UpdateProjectApprovalById(id, projectApproval["ApproverUserPrincipalName"].(string))
			if err != nil {
				logger.LogTrace(err.Error(), contracts.Error)
			}
		}
	}
}

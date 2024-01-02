package routes

import (
	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"
	"net/http"
)

func FillOutApprovers(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	result, err := db.SelectApprovalTypes()
	if err != nil {
		logger.LogException(err)
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
			logger.LogException(err)
		}
	}
}

func FillOutApprovalRequestApprovers(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	projectApprovals, err := db.GetProjectApprovals()
	if err != nil {
		logger.LogException(err)
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
			logger.LogException(err)
		}

		if projectApproval["ApprovalDate"] != nil && projectApproval["ApproverUserPrincipalName"] != nil {
			err = db.UpdateProjectApprovalById(id, projectApproval["ApproverUserPrincipalName"].(string))
			if err != nil {
				logger.LogException(err)
			}
		}
	}
}

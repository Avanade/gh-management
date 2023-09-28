package routes

import (
	"log"
	db "main/pkg/ghmgmtdb"
	"net/http"
)

func FillOutApprovers(w http.ResponseWriter, r *http.Request) {
	result, err := db.SelectApprovalTypes()
	if err != nil {
		log.Println(err.Error())
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
			log.Println(err.Error())
		}
	}
}

func FillOutApprovalRequestApprovers(w http.ResponseWriter, r *http.Request) {
	projectApprovals, err := db.GetProjectApprovals()
	if err != nil {
		log.Println(err.Error())
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
			log.Println(err.Error())
		}

		if projectApproval["ApprovalDate"] != nil {
			err = db.UpdateProjectApprovalById(id, projectApproval["ApproverUserPrincipalName"].(string))
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

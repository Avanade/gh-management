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
		err = db.InsertApprover(db.Approver{
			ApprovalTypeId: int(approver["Id"].(int64)),
			ApproverEmail:  approver["ApproverUserPrincipalName"].(string),
		})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

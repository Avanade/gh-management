package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	"main/pkg/sql"
	"net/http"
	"os"
)

func UpdateApprovalStatusProjects(w http.ResponseWriter, r *http.Request) {

	err := processApprovalProjects(r, "projects")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
func UpdateApprovalStatusCommunity(w http.ResponseWriter, r *http.Request) {
	fmt.Println(1)
	err := processApprovalProjects(r, "community")
	fmt.Println(2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(1)
	w.WriteHeader(http.StatusOK)

}

func processApprovalProjects(r *http.Request, module string) error {

	// Decode payload
	var req models.TypUpdateApprovalStatusReqBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	const REJECTED = 3
	const APPROVED = 5

	// Connect to database
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		return err
	}
	defer db.Close()

	//Update approval status on database
	approvalStatusId := APPROVED
	if !req.IsApproved {
		approvalStatusId = REJECTED
	}

	params := make(map[string]interface{})
	params["ApprovalSystemGUID"] = req.ItemId
	params["ApprovalStatusId"] = approvalStatusId
	params["ApprovalRemarks"] = req.Remarks
	params["ApprovalDate"] = req.ResponseDate

	var spName string
	switch module {
	case "projects":
		spName = "PR_ProjectsApproval_Update_ApproverResponse"
	case "community":
		spName = "PR_CommunityApproval_Update_ApproverResponse"
	}

	_, err = db.ExecuteStoredProcedure(spName, params)
	if err != nil {
		return err
	}
	return nil
}

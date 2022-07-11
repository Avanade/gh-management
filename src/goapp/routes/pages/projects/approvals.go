package routes

import (
	"encoding/json"
	models "main/models"
	ghmgmt "main/pkg/ghmgmtdb"
	gh "main/pkg/github"
	"main/pkg/sql"
	"net/http"
	"os"
	"time"
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
	err := processApprovalProjects(r, "community")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	projectApproval := ghmgmt.GetProjectApprovalByGUID(req.ItemId)

	go checkAllRequests(projectApproval.ProjectId)
	return nil
}

func checkAllRequests(id int64) {
	allApproved := true

	// Check if all requests are approved
	projectApprovals := ghmgmt.GetProjectApprovalsByProjectId(id)
	repo := projectApprovals[0].ProjectName
	for _, a := range projectApprovals {
		if a.RequestStatus != "Approved" {
			allApproved = false
			break
		}
	}

	// If all are approved, move repository to OpenSource and make public
	if allApproved {
		owner := os.Getenv("GH_ORG_INNERSOURCE")
		newOwner := os.Getenv("GH_ORG_OPENSOURCE")
		gh.TransferRepository(repo, owner, newOwner)

		time.Sleep(3 * time.Second)
		gh.SetProjectVisibility(repo, "public", newOwner)

		ghmgmt.UpdateIsPrivate(id, false)
	}
}

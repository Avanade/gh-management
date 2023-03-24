package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	email "main/pkg/email"
	ghmgmt "main/pkg/ghmgmtdb"
	gh "main/pkg/github"
	"main/pkg/sql"
	"net/http"
	"os"
	"strings"
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
	const PUBLIC = 3
	if allApproved {
		owner := os.Getenv("GH_ORG_INNERSOURCE")
		newOwner := os.Getenv("GH_ORG_OPENSOURCE")
		gh.TransferRepository(repo, owner, newOwner)

		time.Sleep(3 * time.Second)
		gh.SetProjectVisibility(repo, "public", newOwner)

		ghmgmt.UpdateProjectVisibilityId(id, PUBLIC)
	}
}

func UpdateApprovalReassignApprover(w http.ResponseWriter, r *http.Request) {
	var req models.TypUpdateApprovalReAssign
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Connect to database
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer db.Close()

	param := map[string]interface{}{
		"Id":            req.Id,
		"ApproverEmail": req.ApproverEmail,
		"Username":      req.Username,
	}

	result, err2 := db.ExecuteStoredProcedureWithResult("PR_ProjectsApproval_Update_ApproverUserPrincipalName", param)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range result {
		data := models.TypProjectApprovals{
			Id:                         v["Id"].(int64),
			ProjectId:                  v["ProjectId"].(int64),
			ProjectName:                v["ProjectName"].(string),
			ProjectCoowner:             v["ProjectCoowner"].(string),
			ProjectDescription:         v["ProjectDescription"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			CoownerGivenName:           v["CoownerGivenName"].(string),
			CoownerSurName:             v["CoownerSurName"].(string),
			CoownerName:                v["CoownerName"].(string),
			CoownerUserPrincipalName:   v["CoownerUserPrincipalName"].(string),
			ApprovalTypeId:             v["ApprovalTypeId"].(int64),
			ApprovalType:               v["ApprovalType"].(string),
			ApproverUserPrincipalName:  v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
			Newcontribution:            v["newcontribution"].(string),
			OSSsponsor:                 v["OSSsponsor"].(string),
			Avanadeofferingsassets:     v["Avanadeofferingsassets"].(string),
			Willbecommercialversion:    v["Willbecommercialversion"].(string),
			OSSContributionInformation: v["OSSContributionInformation"].(string),
			RequestStatus:              v["RequestStatus"].(string),
		}

		err = SendReassignEmail(data)

	}
	w.WriteHeader(http.StatusOK)
}

func SendReassignEmail(data models.TypProjectApprovals) error {

	bodyTemplate := `<p>Hi |ApproverUserPrincipalName|!</p>
		<p>|RequesterName| is requesting for a new project and is now pending for |ApprovalType| review.</p>
		<p>Below are the details:</p>
		<table>
			<tr>
				<td style="font-weight: bold;">Project Name<td>
				<td style="font-size:larger">|ProjectName|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">CoOwner<td>
				<td style="font-size:larger">|CoownerName|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Description<td>
				<td style="font-size:larger">|ProjectDescription|<td>
			</tr>
		</table>
		<table>
			<tr>
				<td style="font-weight: bold;">Is this a new contribution with no prior code development? (i.e., no existing Avanade IP, no third-party/OSS code, etc.)<td>
				<td style="font-size:larger">|Newcontribution|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Who is sponsoring this OSS contribution?<td>
				<td style="font-size:larger">|OSSsponsor|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Will Avanade use this contribution in client accounts and/or as part of an Avanade offerings/assets?<td>
				<td style="font-size:larger">|Avanadeofferingsassets|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Will there be a commercial version of this contribution<td>
				<td style="font-size:larger">|Willbecommercialversion|<td>
			</tr>
				<tr>
				<td style="font-weight: bold;">Additional OSS Contribution Information (e.g. planned maintenance/support, etc.)?<td>
				<td style="font-size:larger">|OSSContributionInformation|<td>
			</tr>
		</table>
		<p>For more information, send an email to <a href="mailto:|RequesterUserPrincipalName|">|RequesterUserPrincipalName|</a></p>
		`
	replacer := strings.NewReplacer("|ApproverUserPrincipalName|", data.ApproverUserPrincipalName,
		"|RequesterName|", data.RequesterName,
		"|ApprovalType|", data.ApprovalType,
		"|ProjectName|", data.ProjectName,
		"|CoownerName|", data.CoownerName,
		"|ProjectDescription|", data.ProjectDescription,
		"|RequesterUserPrincipalName|", data.RequesterUserPrincipalName,

		"|Newcontribution|", data.Newcontribution,
		"|OSSsponsor|", data.OSSsponsor,
		"|Avanadeofferingsassets|", data.Avanadeofferingsassets,
		"|Willbecommercialversion|", data.Willbecommercialversion,
		"|OSSContributionInformation|", data.OSSContributionInformation,
	)

	//buf := new(bytes.Buffer)
	//body := buf.String()
	body := replacer.Replace(bodyTemplate)
	m := email.TypEmailMessage{
		Subject: fmt.Sprintf("[GH-Management] New Project For Review - %v", data.ProjectName),
		Body:    body,
		To:      data.ApproverUserPrincipalName,
	}

	_, errEmail := email.SendEmail(m)

	if errEmail != nil {
		return errEmail
	}
	return errEmail
}

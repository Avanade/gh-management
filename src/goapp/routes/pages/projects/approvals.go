package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"main/pkg/email"
	"main/pkg/envvar"
	db "main/pkg/ghmgmtdb"
	ghAPI "main/pkg/github"
	"main/pkg/notification"
)

type ApprovalReAssignRequestBody struct {
	Id                  string `json:"id"`
	ApproverEmail       string `json:"ApproverEmail"`
	Username            string `json:"Username"`
	ApplicationId       string `json:"ApplicationId"`
	ApplicationModuleId string `json:"ApplicationModuleId"`
	ItemId              string `json:"itemId"`
	ApproveText         string `json:"ApproveText"`
	RejectText          string `json:"RejectText"`
}

type ApprovalStatusRequestBody struct {
	ItemId       string `json:"itemId"`
	IsApproved   bool   `json:"isApproved"`
	Remarks      string `json:"Remarks"`
	ResponseDate string `json:"responseDate"`
}

func UpdateApprovalStatusProjects(w http.ResponseWriter, r *http.Request) {
	err := ProcessApprovalProjects(r, "projects")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateApprovalStatusCommunity(w http.ResponseWriter, r *http.Request) {
	err := ProcessApprovalProjects(r, "community")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateApprovalReassignApprover(w http.ResponseWriter, r *http.Request) {
	var req ApprovalReAssignRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	param := map[string]interface{}{
		"Id":            req.Id,
		"ApproverEmail": req.ApproverEmail,
		"Username":      req.Username,
	}

	result, err := db.ProjectsApprovalUpdateApproverUserPrincipalName(param)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range result {
		data := db.ProjectApproval{
			Id:                         v["Id"].(int64),
			ProjectId:                  v["ProjectId"].(int64),
			ProjectName:                v["ProjectName"].(string),
			ProjectDescription:         v["ProjectDescription"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			ApprovalTypeId:             v["ApprovalTypeId"].(int64),
			ApprovalType:               v["ApprovalType"].(string),
			ApproverUserPrincipalName:  v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
			Newcontribution:            v["newcontribution"].(string),
			OSSsponsor:                 v["OSSsponsor"].(string),
			Offeringsassets:            v["Avanadeofferingsassets"].(string),
			Willbecommercialversion:    v["Willbecommercialversion"].(string),
			OSSContributionInformation: v["OSSContributionInformation"].(string),
			RequestStatus:              v["RequestStatus"].(string),
		}
		data.ApproveUrl = fmt.Sprintf("%s/response/%s/%s/%s/1", os.Getenv("APPROVAL_SYSTEM_APP_URL"), req.ApplicationId, req.ApplicationModuleId, req.ItemId)
		data.RejectUrl = fmt.Sprintf("%s/response/%s/%s/%s/0", os.Getenv("APPROVAL_SYSTEM_APP_URL"), req.ApplicationId, req.ApplicationModuleId, req.ItemId)
		data.ApproveText = req.ApproveText
		data.RejectText = req.RejectText

		err = SendReassignEmail(data)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	w.WriteHeader(http.StatusOK)
}

func UpdateCommunityApprovalReassignApprover(w http.ResponseWriter, r *http.Request) {
	var req ApprovalReAssignRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	param := map[string]interface{}{
		"Id":            req.Id,
		"ApproverEmail": req.ApproverEmail,
		"Username":      req.Username,
	}
	result, err := db.CommunityApprovalslUpdateApproverUserPrincipalName(param)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range result {
		data := db.CommunityApproval{
			Id:                         v["Id"].(int64),
			CommunityId:                v["CommunityId"].(int64),
			CommunityName:              v["ProjectName"].(string),
			CommunityDescription:       v["ProjectDescription"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			CommunityUrl:               v["Url"].(string),
			CommunityNotes:             v["Notes"].(string),
			ApproverUserPrincipalName:  v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
		}
		data.ApproveUrl = fmt.Sprintf("%s/response/%s/%s/%s/1", os.Getenv("APPROVAL_SYSTEM_APP_URL"), req.ApplicationId, req.ApplicationModuleId, req.ItemId)
		data.RejectUrl = fmt.Sprintf("%s/response/%s/%s/%s/0", os.Getenv("APPROVAL_SYSTEM_APP_URL"), req.ApplicationId, req.ApplicationModuleId, req.ItemId)
		data.ApproveText = req.ApproveText
		data.RejectText = req.RejectText

		err = SendReassignEmailCommunity(data)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	w.WriteHeader(http.StatusOK)
}

func SendReassignEmail(data db.ProjectApproval) error {

	bodyTemplate := `<p>Hi |ApproverUserPrincipalName|!</p>
		<p>|RequesterName| is requesting for a new project and is now pending for |ApprovalType| review.</p>
		<p>Below are the details:</p>
		<table>
			<tr>
				<td style="font-weight: bold;">Project Name<td>
				<td style="font-size:larger">|ProjectName|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Requested by<td>
				<td style="font-size:larger">|Requester|<td>
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
		
		<table style="margin: 10px 0;width:100%; text-align: center;">
        <tr>
            <td colspan="5" style="padding: 5px 0;">To process your response, click any of the buttons below:</td>
        </tr>
        
        <tr style="color: white;">
            <td style="padding: 5px 0px; width: 20%; "></td>
            <td style="padding: 5px 0px; width: 26%; background-color: green;">
                <a href="|ApproveUrl|" style="color: white;">
                    |ApproveText|
                </a>
            </td>
            <td style="padding: 5px 0px; width: 8%; "></td>
            <td style="padding: 5px 0px; width: 26%; background-color: red;">
                <a href="|RejectUrl|" style="color: white;">
                    |RejectText|
                </a>
            </td>
            <td style="padding: 5px 0px; width: 20%; "></td>
        </tr>
    </table>
		`
	replacer := strings.NewReplacer("|ApproverUserPrincipalName|", data.ApproverUserPrincipalName,
		"|RequesterName|", data.RequesterName,
		"|ApprovalType|", data.ApprovalType,
		"|ProjectName|", data.ProjectName,
		"|Requester|", data.RequesterName,
		"|ProjectDescription|", data.ProjectDescription,
		"|RequesterUserPrincipalName|", data.RequesterUserPrincipalName,

		"|Newcontribution|", data.Newcontribution,
		"|OSSsponsor|", data.OSSsponsor,
		"|Avanadeofferingsassets|", data.Offeringsassets,
		"|Willbecommercialversion|", data.Willbecommercialversion,
		"|OSSContributionInformation|", data.OSSContributionInformation,
		"|ApproveUrl|", data.ApproveUrl,
		"|RejectUrl|", data.RejectUrl,
		"|ApproveText|", data.ApproveText,
		"|RejectText|", data.RejectText,
	)

	body := replacer.Replace(bodyTemplate)
	m := email.EmailMessage{
		Subject: fmt.Sprintf("[GH-Management] New Project For Review - %v", data.ProjectName),
		Body:    body,
		To:      data.ApproverUserPrincipalName,
	}

	_, err := email.SendEmail(m)

	if err != nil {
		return err
	}
	return nil
}

func ProcessApprovalProjects(r *http.Request, module string) error {

	// Decode payload
	var req ApprovalStatusRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	const REJECTED = 3
	const APPROVED = 5

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

	_, err = db.UpdateApprovalApproverResponse(spName, req.ItemId, req.Remarks, req.ResponseDate, approvalStatusId)
	if err != nil {
		return err
	}

	if module == "projects" {
		projectApproval := db.GetProjectApprovalByGUID(req.ItemId)

		go CheckAllRequests(projectApproval.ProjectId, r.Host)
	}
	return nil
}

func CheckAllRequests(id int64, host string) {
	allApproved := true

	// Check if all requests are approved
	projectApprovals := db.GetProjectApprovalsByProjectId(id)
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
		ghAPI.TransferRepository(repo, owner, newOwner)

		time.Sleep(3 * time.Second)
		ghAPI.SetProjectVisibility(repo, "public", newOwner)

		db.UpdateProjectVisibilityId(id, PUBLIC)

		repoResp, _ := ghAPI.GetRepository(repo, newOwner)
		db.UpdateTFSProjectReferenceById(id, repoResp.GetHTMLURL())

		repoOwners, err := db.GetRepoOwnersRecordByRepoId(id)
		if err != nil {
			log.Println(err.Error())
			return
		}

		var recipients []string

		for i, _ := range repoOwners {
			recipients = append(recipients, repoOwners[i].UserPrincipalName)
		}

		messageBody := notification.RepositoryPublicApprovalProvidedMessageBody{
			Recipients:          recipients,
			CommunityPortalLink: fmt.Sprint(envvar.GetEnvVar("SCHEME", "https"), "://", host, "/repositories"),
			RepoLink:            repoResp.GetHTMLURL(),
			RepoName:            repoResp.GetName(),
		}
		err = messageBody.Send()
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func SendReassignEmailCommunity(data db.CommunityApproval) error {

	bodyTemplate := `<p>Hi |ApproverUserPrincipalName|!</p>
		<p>|RequesterName| is requesting for a new |CommunityType| community and is now pending for approval.</p>
		<p>Below are the details:</p>
		<table>
			<tr>
				<td style="font-weight: bold;">Community Name<td>
				<td style="font-size:larger">|CommunityName|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Url<td>
				<td style="font-size:larger">|CommunityUrl|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Description<td>
				<td style="font-size:larger">|CommunityDescription|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Notes<td>
				<td style="font-size:larger">|CommunityNotes|<td>
			</tr>
		</table>
		<p>For more information, send an email to <a href="mailto:|RequesterUserPrincipalName|">|RequesterUserPrincipalName|</a></p>

		<table style="margin: 10px 0;width:100%; text-align: center;">
        <tr>
            <td colspan="5" style="padding: 5px 0;">To process your response, click any of the buttons below:</td>
        </tr>
        
        <tr style="color: white;">
            <td style="padding: 5px 0px; width: 20%; "></td>
            <td style="padding: 5px 0px; width: 26%; background-color: green;">
                <a href="|ApproveUrl|" style="color: white;">
                    |ApproveText|
                </a>
            </td>
            <td style="padding: 5px 0px; width: 8%; "></td>
            <td style="padding: 5px 0px; width: 26%; background-color: red;">
                <a href="|RejectUrl|" style="color: white;">
                    |RejectText|
                </a>
            </td>
            <td style="padding: 5px 0px; width: 20%; "></td>
        </tr>
    </table>
		`

	replacer := strings.NewReplacer("|ApproverUserPrincipalName|", data.ApproverUserPrincipalName,
		"|RequesterName|", data.RequesterName,
		"|CommunityType|", data.CommunityType,
		"|CommunityName|", data.CommunityName,
		"|CommunityUrl|", data.CommunityUrl,
		"|CommunityDescription|", data.CommunityDescription,
		"|CommunityNotes|", data.CommunityNotes,
		"|RequesterUserPrincipalName|", data.RequesterUserPrincipalName,
		"|ApproveUrl|", data.ApproveUrl,
		"|RejectUrl|", data.RejectUrl,
		"|ApproveText|", data.ApproveText,
		"|RejectText|", data.RejectText,
	)
	body := replacer.Replace(bodyTemplate)
	m := email.EmailMessage{
		Subject: fmt.Sprintf("[GH-Management] New Community For Approval - %v", data.CommunityName),
		Body:    body,
		To:      data.ApproverUserPrincipalName,
	}

	_, err := email.SendEmail(m)

	if err != nil {
		return err
	}
	return nil
}

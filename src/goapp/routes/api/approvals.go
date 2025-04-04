package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"main/pkg/appinsights_wrapper"
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
	ApproveText         string `json:"ApproveText"`
	RejectText          string `json:"RejectText"`
}

type ApprovalStatusRequestBody struct {
	ItemId       string `json:"itemId"`
	IsApproved   bool   `json:"isApproved"`
	Remarks      string `json:"Remarks"`
	ResponseDate string `json:"responseDate"`
	RespondedBy  string `json:"respondedBy"`
}

func UpdateApprovalStatusProjects(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	err := ProcessApprovalProjects(r, "projects")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateApprovalStatusCommunity(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	err := ProcessApprovalProjects(r, "community")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateApprovalStatusOrganization(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	err := ProcessApprovalProjects(r, "organization")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateApprovalStatusCopilot(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	err := ProcessApprovalProjects(r, "github-copilot")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateApprovalStatusOrganizationAccess(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	err := ProcessApprovalProjects(r, "orgaccess")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateApprovalStatusAdoOrganization(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	err := ProcessApprovalProjects(r, "ado")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateApprovalReassignApprover(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	var req ApprovalReAssignRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.LogException(err)
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
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range result {
		data := db.ProjectApproval{
			Id:                         v["Id"].(int64),
			ProjectId:                  v["RepositoryId"].(int64),
			ProjectName:                v["RepositoryName"].(string),
			ProjectDescription:         v["RepositoryDescription"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			ApprovalTypeId:             v["RepositoryApprovalTypeId"].(int64),
			ApprovalType:               v["ApprovalType"].(string),
			ApproverUserPrincipalName:  v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
			Newcontribution:            v["Newcontribution"].(string),
			OSSsponsor:                 v["OSSsponsor"].(string),
			Offeringsassets:            v["Avanadeofferingsassets"].(string),
			Willbecommercialversion:    v["Willbecommercialversion"].(string),
			OSSContributionInformation: v["OSSContributionInformation"].(string),
			RequestStatus:              v["RequestStatus"].(string),
		}
		data.ApproveUrl = fmt.Sprintf("%s/response/%s/%s/%s/1", os.Getenv("APPROVAL_SYSTEM_APP_URL"), req.ApplicationId, req.ApplicationModuleId, req.Id)
		data.RejectUrl = fmt.Sprintf("%s/response/%s/%s/%s/0", os.Getenv("APPROVAL_SYSTEM_APP_URL"), req.ApplicationId, req.ApplicationModuleId, req.Id)
		data.ApproveText = req.ApproveText
		data.RejectText = req.RejectText

		err = SendReassignEmail(data)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	w.WriteHeader(http.StatusOK)
}

func UpdateCommunityApprovalReassignApprover(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	var req ApprovalReAssignRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	param := map[string]interface{}{
		"ApprovalSystemGUID":        req.Id,
		"ApproverUserPrincipalName": req.ApproverEmail,
		"UserPrincipalName":         req.Username,
	}
	result, err := db.CommunityApprovalslUpdateApproverUserPrincipalName(param)
	if err != nil {
		logger.LogException(err)
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
		data.ApproveUrl = fmt.Sprintf("%s/response/%s/%s/%s/1", os.Getenv("APPROVAL_SYSTEM_APP_URL"), req.ApplicationId, req.ApplicationModuleId, req.Id)
		data.RejectUrl = fmt.Sprintf("%s/response/%s/%s/%s/0", os.Getenv("APPROVAL_SYSTEM_APP_URL"), req.ApplicationId, req.ApplicationModuleId, req.Id)
		data.ApproveText = req.ApproveText
		data.RejectText = req.RejectText

		err = SendReassignEmailCommunity(data)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
	w.WriteHeader(http.StatusOK)
}

func SendReassignEmail(data db.ProjectApproval) error {

	bodyTemplate := `
		<html>
			<head>
				<style>
					table,
					th,
					tr,
					td {
					border: 0;
					border-collapse: collapse;
					vertical-align: middle;
					}

					.thead {
					padding: 15px;
					}

					.center-table {
					text-align: -webkit-center;
					}

					.margin-auto {
					margin: auto;
					}

					.border-top {
					border-top: 1px rgb(204, 204, 204) solid;
					border-collapse: separate;
					}
				</style>
			</head>

			<body>
				<table style="width: 100%">
					<tr>
						<th class="center-table">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr style="background-color: #ff5800">
									<td class="thead" style="width: 95px">
									<img
										src="data:image/png;base64,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" />
									</td>
									<td class="thead" style="
										font-family: SegoeUI, sans-serif;
										font-size: 14px;
										color: white;
										padding-bottom: 10px;
										">
									Community
									</td>
								</tr>
							</table>
						</th>
					</tr>
					<tr>
						<th class="center-table">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="font-size: 18px; font-weight: 600; padding-top: 20px">
									Request for |ApprovalType| Review
									</td>
								</tr>
							</table>
						</th>
					</tr>
					<tr>
						<td class="center-table"  align="center">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="padding-top: 20px">
										<p>Hi,</p>
										<p>|RequesterName| is requesting for a repository to be public and is now pending for |ApprovalType| review.</p>
										<p>Below are the details:</p>
									</td>
								</tr>
							</table>
						</td>
					</tr>
					<tr>
						<th class="center-table">
							<table style="width: 100%; max-width: 700px; margin: auto">
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Repository Name
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|ProjectName|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Requested by
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Requester|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Description
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|ProjectDescription|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Is this a new contribution with no prior code development? <br>
										(i.e., no existing |OrganizationName| IP, no third-party/OSS code, etc.)
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Newcontribution|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Who is sponsoring this OSS contribution?
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|OSSsponsor|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Will |OrganizationName| use this contribution in client accounts <br>
										and/or as part of an |OrganizationName| offerings/assets?
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Avanadeofferingsassets|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Will there be a commercial version of this contribution?
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Willbecommercialversion|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Additional OSS Contribution Information <br>
										(e.g. planned maintenance/support, etc.)?
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|OSSContributionInformation|
									</td>
								</tr>
							</table>
						</th>
					</tr>
					<tr>
						<td class="center-table">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="padding-top: 20px">
									For more information, send an email to <a href="mailto:|RequesterUserPrincipalName|">|RequesterUserPrincipalName|</a>
									</td>
								</tr>
							</table>
						</td>
					</tr>
					<tr>
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
					</tr>
				</table>
				<br>
			</body>

		</html>
		`

	replacer := strings.NewReplacer(
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
		"|OrganizationName|", os.Getenv("ORGANIZATION_NAME"),
	)

	body := replacer.Replace(bodyTemplate)
	m := email.Message{
		Subject: fmt.Sprintf("Request for %v Review - %v", data.ApprovalType, data.ProjectName),
		Body: email.Body{
			Content: body,
			Type:    email.HtmlMessageType,
		},
		ToRecipients: []email.Recipient{
			{
				Email: data.ApproverUserPrincipalName,
			},
		},
	}

	err := email.SendEmail(m, true)

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

	switch module {
	case "projects":
		_, err = db.UpdateProjectApprovalApproverResponse(req.ItemId, req.Remarks, req.ResponseDate, req.RespondedBy, approvalStatusId)
		if err != nil {
			return err
		}

		projectApproval := db.GetProjectApprovalByGUID(req.ItemId)

		go CheckAllRequests(projectApproval.ProjectId, r.Host)
	case "community":
		_, err = db.UpdateCommunityApprovalApproverResponse(req.ItemId, req.Remarks, req.ResponseDate, approvalStatusId)
		if err != nil {
			return err
		}
	case "organization":
		_, err = db.UpdateApprovalApproverResponse(req.ItemId, req.Remarks, req.ResponseDate, approvalStatusId, req.RespondedBy)
		if err != nil {
			return err
		}
	case "github-copilot":
		_, err = db.UpdateApprovalApproverResponse(req.ItemId, req.Remarks, req.ResponseDate, approvalStatusId, req.RespondedBy)
		if err != nil {
			return err
		}

		// If approved add the user to GitHub Copilot License Group
		if approvalStatusId == APPROVED {
			gc, err := db.GetGitHubCopilotbyGUID(req.ItemId)
			if err != nil {
				return err
			}
			org := gc[0]["RegionName"].(string)
			ghUsername := gc[0]["GitHubUsername"].(string)

			// Check if team is existing
			ghToken := os.Getenv("GH_TOKEN")
			slug := os.Getenv("COPILOT_GROUP_SLUG")
			team, err := ghAPI.GetTeam(ghToken, org, slug)
			if err != nil {
				return err
			}

			// If not existing create the team
			if team == nil {
				_, err := ghAPI.CreateTeam(ghToken, org, "GitHub Copilot License Group")
				if err != nil {
					return err
				}
			}

			// Add user to the team
			_, err = ghAPI.AddMemberToTeam(ghToken, org, slug, ghUsername, "member")
			if err != nil {
				return err
			}
		}
	case "orgaccess":
		_, err = db.UpdateApprovalApproverResponse(req.ItemId, req.Remarks, req.ResponseDate, approvalStatusId, req.RespondedBy)
		if err != nil {
			return err
		}

		if approvalStatusId == APPROVED {
			orgAccess, err := db.GetOrganizationAccessByApprovalRequestItemId(req.ItemId)
			if err != nil {
				return err
			}
			ghAPI.OrganizationInvitation(os.Getenv("GH_TOKEN"), orgAccess.User.GithubUsername, orgAccess.Organization.Name)
		}
	case "ado":
		_, err = db.UpdateApprovalApproverResponse(req.ItemId, req.Remarks, req.ResponseDate, approvalStatusId, req.RespondedBy)
		if err != nil {
			return err
		}

		if approvalStatusId == APPROVED {
			// TODO: Trigger function app to create ADO Organization
		}
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

		ValidateOrgMembers(owner, repo, newOwner, nil)
		ghAPI.SetProjectVisibility(repo, "public", owner)
		ghAPI.TransferRepository(repo, owner, newOwner)
		time.Sleep(3 * time.Second)
		db.UpdateProjectVisibilityId(id, PUBLIC)

		repoResp, _ := ghAPI.GetRepository(repo, newOwner)
		db.UpdateTFSProjectReferenceById(id, repoResp.GetHTMLURL(), *repoResp.GetOwner().Login)
	}

	// Check if all requests are responded by approvers.
	allResponded := true

	for _, a := range projectApprovals {
		if a.ApprovalDate.IsZero() {
			allResponded = false
			break
		}
	}

	// Send notification if all approvers responded to the approval request
	if allResponded {
		repoOwners, err := db.GetRepoOwnersRecordByRepoId(id)
		if err != nil {
			log.Println(err.Error())
			return
		}

		var recipients []string

		for i := range repoOwners {
			recipients = append(recipients, repoOwners[i].UserPrincipalName)
		}

		project := db.GetProjectById(id)

		messageBody := notification.RepositoryPublicApprovalProvidedMessageBody{
			Recipients:          recipients,
			CommunityPortalLink: fmt.Sprint(envvar.GetEnvVar("SCHEME", "https"), "://", host, "/repositories"),
			RepoLink:            project[0]["TFSProjectReference"].(string),
			RepoName:            project[0]["Name"].(string),
		}
		err = messageBody.Send()
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func SendReassignEmailCommunity(data db.CommunityApproval) error {

	bodyTemplate := `
		<html>
			<head>
				<style>
					table,
					th,
					tr,
					td {
					border: 0;
					border-collapse: collapse;
					vertical-align: middle;
					}

					.thead {
					padding: 15px;
					}

					.center-table {
					text-align: -webkit-center;
					}

					.margin-auto {
					margin: auto;
					}

					.border-top {
					border-top: 1px rgb(204, 204, 204) solid;
					border-collapse: separate;
					}
				</style>
			</head>

			<body>
				<table style="width: 100%">
					<tr>
						<th class="center-table">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr style="background-color: #ff5800">
									<td class="thead" style="width: 95px">
									<img
										src="data:image/png;base64,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" />
									</td>
									<td class="thead" style="
										font-family: SegoeUI, sans-serif;
										font-size: 14px;
										color: white;
										padding-bottom: 10px;
										">
									Community
									</td>
								</tr>
							</table>
						</th>
					</tr>
					<tr>
						<th class="center-table">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="font-size: 18px; font-weight: 600; padding-top: 20px">
									Request for a New Community
									</td>
								</tr>
							</table>
						</th>
					</tr>
					<tr>
						<td class="center-table"  align="center">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="padding-top: 20px">
										<p>Hi |ApproverUserPrincipalName|!</p>
										<p>|RequesterName| is requesting for a new |CommunityType| community and is now pending for approval.</p>
										<p>Below are the details:</p>
									</td>
								</tr>
							</table>
						</td>
					</tr>
					<tr>
						<th class="center-table">
							<table style="width: 100%; max-width: 700px; margin: auto">
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Community Name
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|CommunityName|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Url
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|CommunityUrl|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Description
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|CommunityDescription|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Notes
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|CommunityNotes|
									</td>
								</tr>
							</table>
						</th>
					</tr>
					<tr>
						<td class="center-table">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="padding-top: 20px">
									For more information, send an email to <a href="mailto:|RequesterUserPrincipalName|">|RequesterUserPrincipalName|</a>
									</td>
								</tr>
							</table>
						</td>
					</tr>
					<tr>
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
					</tr>
				</table>
				<br>
			</body>

		</html>
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
	m := email.Message{
		Subject: fmt.Sprintf("New Community For Approval - %v", data.CommunityName),
		Body: email.Body{
			Content: body,
			Type:    email.HtmlMessageType,
		},
		ToRecipients: []email.Recipient{
			{
				Email: data.ApproverUserPrincipalName,
			},
		},
	}

	err := email.SendEmail(m, true)

	if err != nil {
		return err
	}
	return nil
}

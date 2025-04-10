package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"
	ghAPI "main/pkg/github"
	"main/pkg/notification"
	"main/pkg/session"

	"github.com/gorilla/mux"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
)

// GITHUB COPILOT LICENSE
func AddGitHubCopilot(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	// Get username
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	// Get GitHub ID
	sessiongh, _ := session.Store.Get(r, "gh-auth-session")
	ghProfile := sessiongh.Values["ghProfile"].(string)
	var p map[string]interface{}
	err := json.Unmarshal([]byte(ghProfile), &p)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ghId := p["id"].(float64)
	ghUser := fmt.Sprintf("%s", p["login"])

	// Parse request body
	var body db.GitHubCopilot
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body.Username = username.(string)
	body.GitHubId = int64(ghId)
	body.GitHubUsername = ghUser

	// Check user's membership
	token := os.Getenv("GH_TOKEN")

	// Check if there is a pending request
	result, err := db.GitHubCopilotGetPendingByUserAndOrganization(body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result != nil {
		http.Error(w, "You have a pending request on this organization. Wait for response from the approvers. ", http.StatusBadRequest)
		return
	}

	// Get organization owners to produce list of approvers
	var approvers []string
	orgOwners, err := ghAPI.OrgListMembers(token, body.RegionName, "admin")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exemption := os.Getenv("EXEMPTION")

	for _, owner := range orgOwners {
		email, err := db.GetUserEmailByGithubId(strconv.FormatInt(*owner.ID, 10))
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if email != "" && email != exemption {
			approvers = append(approvers, email)
		}
	}

	if len(approvers) == 0 {
		logger.LogException("Can't find email address of organization owners.")
		http.Error(w, "Can't find email address of organization owners.", http.StatusInternalServerError)
		return
	}

	// Insert record on GitHubCopilot table
	result, err = db.GitHubCopilotInsert(body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var requestIds []int64
	for _, approver := range approvers {

		// Insert approval request record
		requestId, err := db.ApprovalInsert(approver, "GitHub Copilot License Request", body.Username)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		requestIds = append(requestIds, requestId)

		// Insert link record
		err = db.GitHubCopilotApprovalInsert(id, requestId)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	body.Approvers = approvers
	body.RequestIds = requestIds
	body.Id = int64(id)

	CreateGitHubCopilotApprovalRequest(body, logger)

	messageBody := notification.RequestForGitHubCopilotLicenseMessageBody{
		Recipients: approvers,
		UserName:   ghUser,
	}

	err = messageBody.Send()
	if err != nil {
		logger.LogException(err)
	}
}

func CreateGitHubCopilotApprovalRequest(data db.GitHubCopilot, logger *appinsights_wrapper.TelemetryClient) error {

	url := os.Getenv("APPROVAL_SYSTEM_APP_URL")
	if url != "" {
		url = url + "/api/request"

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
									Request for a GitHub Copilot License
									</td>
								</tr>
							</table>
						</th>
					</tr>
					<tr>
						<th class="center-table">
							<table style="width: 100%; max-width: 700px; margin: auto">
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										GitHub Organization
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Region|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										GitHub Id
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|GitHubId|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										GitHub Username
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|GitHubUsername|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Requested By
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|RequestedBy|
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
										Note: Approving the request will add the requestor to GitHub Copilot License Group team on |Region|
									</td>
								</tr>
							</table>
						</td>
					</tr>
				</table>
				<br>
			</body>

		</html>
		`

		replacer := strings.NewReplacer("|Region|", data.RegionName,
			"|GitHubId|", strconv.FormatInt(data.GitHubId, 10),
			"|GitHubUsername|", data.GitHubUsername,
			"|RequestedBy|", data.Username,
		)
		body := replacer.Replace(bodyTemplate)

		postParams := CommunityApprovalSystemPost{
			ApplicationId:       os.Getenv("APPROVAL_SYSTEM_APP_ID"),
			ApplicationModuleId: os.Getenv("APPROVAL_SYSTEM_APP_MODULE_COPILOT"),
			Emails:              data.Approvers,
			Subject:             "New GitHub Copilot License Request",
			Body:                body,
			RequesterEmail:      data.Username,
		}

		r := getHttpPostResponseStatus(url, postParams, logger)
		if r != nil {
			var res CommunityApprovalSystemPostResponseDto
			err := json.NewDecoder(r.Body).Decode(&res)
			if err != nil {
				return err
			}
			for _, reqId := range data.RequestIds {
				db.CommunityApprovalUpdateGUID(reqId, res.ItemId)
			}
		}
	}
	return nil
}

func GetAllGitHubCopilotRequest(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	// Get username
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	orgs, err := db.GetAllGitHubCopilotRequest(username.(string))
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(orgs)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetGitHubCopilotApprovalRequests(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id, err := strconv.ParseInt(req["id"], 10, 64)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	approvals, err := db.GetGitHubCopilotApprovalRequest(id)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvals)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func ReprocessCommunityApprovalRequestGitHubCoPilots() {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	items := db.GetFailedCommunityApprovalRequestGitHubCoPilots()

	for _, item := range items {
		err := CreateGitHubCopilotApprovalRequest(item, logger)
		if err != nil {
			logger.LogTrace("ID:"+strconv.FormatInt(item.Id, 10)+" "+err.Error(), contracts.Error)
		}
	}
}

// ORGANIZATION ACCESS
func RequestOrganizationAccess(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	// Get username
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"].(string)

	// Get GitHub ID
	sessiongh, _ := session.Store.Get(r, "gh-auth-session")
	ghProfile := sessiongh.Values["ghProfile"].(string)
	var p map[string]interface{}
	err := json.Unmarshal([]byte(ghProfile), &p)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ghUsername := fmt.Sprintf("%s", p["login"])

	// Parse request body
	var regionalOrg struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}

	err = json.NewDecoder(r.Body).Decode(&regionalOrg)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token := os.Getenv("GH_TOKEN")
	membership, _ := ghAPI.UserMembership(token, regionalOrg.Name, ghUsername)
	if membership != nil {
		switch membership.GetState() {
		case "active":
			logger.LogException(err)
			http.Error(w, "The request cannot proceed because you are already a member of this organization.", http.StatusBadRequest)
			return
		case "pending":
			logger.LogException(err)
			http.Error(w, fmt.Sprint("The request cannot proceed because you have pending invitation from this organization, ", regionalOrg.Name), http.StatusBadRequest)
			return
		}
	}

	hasPendingRequest, err := db.HasOrganizationAccessPendingRequest(username, regionalOrg.Id)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if hasPendingRequest {
		logger.LogException(err)
		http.Error(w, "The request cannot proceed because you have pending request.", http.StatusBadRequest)
		return
	}

	// Get organization owners to produce list of approvers
	var approvers []string
	orgOwners, err := ghAPI.OrgListMembers(token, regionalOrg.Name, "admin")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exemption := os.Getenv("EXEMPTION")

	for _, owner := range orgOwners {
		email, err := db.GetUserEmailByGithubId(strconv.FormatInt(*owner.ID, 10))
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if email != "" && email != exemption {
			approvers = append(approvers, email)
		}
	}

	if len(approvers) == 0 {
		logger.LogException("Can't find email address of organization owners.")
		http.Error(w, "Can't find email address of organization owners.", http.StatusInternalServerError)
		return
	}

	// Insert record on OrganizationAccess table
	id, err := db.InsertOrganizationAccess(username, regionalOrg.Id)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var requestIds []int64
	for _, approver := range approvers {
		// Insert approval request record
		requestId, err := db.ApprovalInsert(approver, "Organization Access Request", username)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		requestIds = append(requestIds, requestId)

		// Insert link record
		err = db.InsertOrganizationAccessApprovalRequest(id, requestId)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	CreateOrganizationAccessApprovalRequest(regionalOrg.Name, ghUsername, username, approvers, requestIds, logger)
}

func GetMyOrganizationAccess(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	// Get username
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"].(string)

	organizationAccess, err := db.GetOrganizationAccessByUserPrincipalName(username)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(organizationAccess)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func CreateOrganizationAccessApprovalRequest(
	regionName, ghUsername, username string,
	approverUserPrincipalName []string,
	requestIds []int64,
	logger *appinsights_wrapper.TelemetryClient,
) error {
	url := os.Getenv("APPROVAL_SYSTEM_APP_URL")
	if url != "" {
		url = url + "/api/request"

		bodyTemplate := `<html>
							<head>
								<style>
									table, th, tr, td {
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
														Request for Organization Access
													</td>
												</tr>
											</table>
										</th>
									</tr>
									<tr>
										<th class="center-table">
											<table style="width: 100%; max-width: 700px; margin: auto">
												<tr class="border-top">
													<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
														GitHub Organization
													</td>
													<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
														|Region|
													</td>
												</tr>
												<tr class="border-top">
													<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
														GitHub Username
													</td>
													<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
														|GitHubUsername|
													</td>
												</tr>
												<tr class="border-top">
													<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
														Requested By
													</td>
													<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
														|RequestedBy|
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
													Note: Approving the request will create an invitation for the requestor to join the organization.
													</td>
												</tr>
											</table>
										</td>
									</tr>
								</table>
							</body>
						</html>`

		replacer := strings.NewReplacer(
			"|Region|", regionName,
			"|GitHubUsername|", ghUsername,
			"|RequestedBy|", username,
		)
		body := replacer.Replace(bodyTemplate)

		postParams := CommunityApprovalSystemPost{
			ApplicationId:       os.Getenv("APPROVAL_SYSTEM_APP_ID"),
			ApplicationModuleId: os.Getenv("APPROVAL_SYSTEM_APP_MODULE_ORGACCESS"),
			Emails:              approverUserPrincipalName,
			Subject:             "New Organization Access Request",
			Body:                body,
			RequesterEmail:      username,
		}

		r := getHttpPostResponseStatus(url, postParams, logger)
		if r != nil {
			var res CommunityApprovalSystemPostResponseDto
			err := json.NewDecoder(r.Body).Decode(&res)
			if err != nil {
				return err
			}

			url := fmt.Sprintf("%v/response/%v/%v/%v/1", os.Getenv("APPROVAL_SYSTEM_APP_URL"), os.Getenv("APPROVAL_SYSTEM_APP_ID"), os.Getenv("APPROVAL_SYSTEM_APP_MODULE_ORGACCESS"), res.ItemId)

			messageBody := notification.RequestForOrganizationAccessMessageBody{
				Recipients:       approverUserPrincipalName,
				OrganizationName: regionName,
				OrganizationLink: fmt.Sprintf("https://github.com/%v", regionName),
				ApprovalLink:     url,
				UserName:         ghUsername,
			}

			err = messageBody.Send()
			if err != nil {
				logger.LogException(err)
			}

			for _, requestId := range requestIds {
				db.CommunityApprovalUpdateGUID(requestId, res.ItemId)
			}
		}
	}
	return nil
}

func GetOrganizationAccessApprovalRequests(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id, err := strconv.ParseInt(req["id"], 10, 64)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	approvals, err := db.GetOrganizationAccessApprovalRequest(id)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvals)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func ReprocessCommunityApprovalRequestOrganizationAccess() {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	items := db.GetFailedCommunityApprovalRequestOrganizationAccess()

	for _, item := range items {
		err := CreateOrganizationAccessApprovalRequest(
			item.RegionName,
			item.GitHubUsername,
			item.UserPrincipalName,
			item.Approvers,
			item.RequestIds,
			logger,
		)
		if err != nil {
			logger.LogTrace("ID:"+strconv.FormatInt(item.Id, 10)+" "+err.Error(), contracts.Error)
		}
	}
}

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
	"main/pkg/session"

	"github.com/gorilla/mux"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
)

func AddOrganization(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	// Get username
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	var body db.OrganizationDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body.Username = username.(string)

	// Insert record on organization table
	result, err := db.OrganizationInsert(body)
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

	// Get approver list
	approvers, err := db.GetActiveCommunityApprovers("organization")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var approverList []string
	var requestIds []int64
	for _, approver := range approvers {
		approverUserPrincipalName := approver["ApproverUserPrincipalName"].(string)
		approverList = append(approverList, approverUserPrincipalName)

		// Insert approval request record
		requestId, err := db.ApprovalInsert(approverUserPrincipalName, fmt.Sprintf("GitHub Organization for %s - %s", body.ClientName, body.ProjectName), body.Username)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		requestIds = append(requestIds, requestId)

		// Insert link record
		err = db.OrganizationApprovalInsert(id, requestId)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	// Create approval request
	regOrg, err := db.GetRegionalOrganizationById(body.Region)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body.RegionName = regOrg[0]["Name"].(string)
	body.ApproverUserPrincipalName = approverList
	body.RequestId = requestIds
	body.Id = int64(id)
	err = CreateOrganizationApprovalRequest(body, logger)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateOrganizationApprovalRequest(data db.OrganizationDto, logger *appinsights_wrapper.TelemetryClient) error {

	url := os.Getenv("APPROVAL_SYSTEM_APP_URL")
	if url != "" {
		url = url + "/api/request"
		ch := make(chan *http.Response)

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
									Request for a GitHub Enterprise Organization
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
										Region
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Region|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Client Name
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|ClientName|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Project Name
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|ProjectName|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										WBS
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|WBS|
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
				</table>
			</body>

		</html>
		`

		replacer := strings.NewReplacer("|Region|", data.RegionName,
			"|ClientName|", data.ClientName,
			"|ProjectName|", data.ProjectName,
			"|WBS|", data.WBS,
			"|RequestedBy|", data.Username,
		)
		body := replacer.Replace(bodyTemplate)

		postParams := CommunityApprovalSystemPost{
			ApplicationId:       os.Getenv("APPROVAL_SYSTEM_APP_ID"),
			ApplicationModuleId: os.Getenv("APPROVAL_SYSTEM_APP_MODULE_ORGANIZATION"),
			Emails:              data.ApproverUserPrincipalName,
			Subject:             fmt.Sprintf("[GH-Management] New Organization Request - %v - %v", data.ClientName, data.ProjectName),
			Body:                body,
			RequesterEmail:      data.Username,
		}

		go getHttpPostResponseStatus(url, postParams, ch, logger)
		r := <-ch
		if r != nil {
			var res CommunityApprovalSystemPostResponseDto
			err := json.NewDecoder(r.Body).Decode(&res)
			if err != nil {
				return err
			}
			for _, reqId := range data.RequestId {
				db.CommunityApprovalUpdateGUID(reqId, res.ItemId)
			}
			db.OrganizationUpdateApprovalStatus(data.Id, 2)
		}
	}
	return nil
}

func GetAllRegionalOrganizations(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	regOrgs, err := db.GetAllRegionalOrganizations()
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(regOrgs)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetAllActiveOrganizationApprovers(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	approvers, err := db.GetActiveCommunityApprovers("organization")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvers)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetAllOrganizationRequest(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	// Get username
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	orgs, err := db.GetAllOrganizationRequest(username.(string))
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

func GetOrganizationApprovalRequests(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id, err := strconv.ParseInt(req["id"], 10, 64)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	approvals, err := db.GetOrganizationApprovalRequest(id)
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

func IndexRegionalOrganizations(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	token := os.Getenv("GH_TOKEN")
	orgs, err := ghAPI.GetOrganizations(token)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, org := range orgs {
		prefix := os.Getenv("REGIONAL_ORG_PREFIX")
		if strings.HasPrefix(strings.ToLower(*org.Login), prefix) {
			err = db.RegionalOrganizationInsert(*org.ID, *org.Login)
			logger.LogTrace("Indexing "+*org.Login, contracts.Information)
			if err != nil {
				logger.LogException(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

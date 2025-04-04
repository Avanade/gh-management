package adoOrganization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/config"
	"main/model"
	"main/pkg/appinsights_wrapper"
	"main/pkg/authentication"
	"main/pkg/session"
	"main/service"
	"net/http"
	"strings"
	"time"
)

type adoOrganizationController struct {
	Service *service.Service
	Conf    config.ConfigManager
}

func NewAdoOrganizationController(service *service.Service, conf config.ConfigManager) AdoOrganizationController {
	return &adoOrganizationController{service, conf}
}

func (c *adoOrganizationController) CreateAdoOrganizationRequest(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	// Get request body
	var adoOrganizationRequest model.AdoOrganizationRequest
	err := json.NewDecoder(r.Body).Decode(&adoOrganizationRequest)
	if err != nil {
		logger.TrackException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user details
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	adoOrganizationRequest.CreatedBy = username

	// Insert ADO organization request record
	id, err := c.Service.AdoOrganization.Insert(&adoOrganizationRequest)
	if err != nil {
		logger.TrackException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get approvers
	approvers, err := c.Service.CommunityApprover.GetByCategory("ado-organization")
	if err != nil {
		logger.TrackException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a list of approver user principal names from approvers
	approversUserPrincipalNames := []string{}
	for _, approver := range approvers {
		approversUserPrincipalNames = append(approversUserPrincipalNames, approver.ApproverUserPrincipalName)
	}

	var requestIds []int64

	// Loop through approvers
	for _, approver := range approvers {
		// Insert approval request record
		requestId, err := c.Service.ApprovalRequest.Insert(approver.ApproverUserPrincipalName, fmt.Sprintf("New ADO Organization Request - %s", adoOrganizationRequest.Name), username)
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		requestIds = append(requestIds, requestId)

		// Insert link record between ADO organization and approval request
		err = c.Service.AdoOrganizationApprovalRequest.Insert(id, requestId)
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Send create request to approval system
	url := c.Conf.GetApprovalSystemAppUrl()
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
										Request for ADO Organization
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
											Name
										</td>
										<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
											|Name|
										</td>
									</tr>
									<tr class="border-top">
										<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
											Purpose
										</td>
										<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
											|Purpose|
										</td>
									</tr>
									<tr class="border-top">
										<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
											Requested by
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
											Note: Approving the request will automatically create the ADO organization.
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

		replacer := strings.NewReplacer(
			"|Name|", adoOrganizationRequest.Name,
			"|Purpose|", adoOrganizationRequest.Purpose,
			"|RequestedBy|", adoOrganizationRequest.CreatedBy,
		)

		body := replacer.Replace(bodyTemplate)

		postParams := model.ApprovalSystemRequestBody{
			ApplicationId:       c.Conf.GetApprovalSystemAppId(),
			ApplicationModuleId: c.Conf.GetApprovalSystemAppModuleAdoOrganization(),
			Emails:              approversUserPrincipalNames,
			Subject:             fmt.Sprintf("New ADO Organization Request - %v", adoOrganizationRequest.Name),
			Body:                body,
			RequesterEmail:      adoOrganizationRequest.CreatedBy,
		}

		jsonReq, err := json.Marshal(postParams)
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Generate token
		token, err := authentication.GenerateToken()
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{
			Timeout: time.Second * 30,
		}
		res, err := client.Do(req)
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer res.Body.Close()

		var result model.ApprovalSystemResponse
		err = json.NewDecoder(res.Body).Decode(&result)
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, reqId := range requestIds {
			// Save approval system request id
			err = c.Service.ApprovalRequest.UpdateApprovalSystemGUID(reqId, result.ItemId)
			if err != nil {
				logger.TrackException(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	// Send http status ok
	w.WriteHeader(http.StatusOK)

}

func (c *adoOrganizationController) GetAdoOrganizationByUser(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	adoOrgRequests, err := c.Service.AdoOrganization.GetByUser(username)
	if err != nil {
		logger.TrackException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(adoOrgRequests)
}

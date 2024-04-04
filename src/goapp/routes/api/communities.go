package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"

	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"
	"main/pkg/msgraph"
	"main/pkg/session"

	"github.com/gorilla/mux"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
	"github.com/thedatashed/xlsxreader"
)

type CommunityDto struct {
	Id                     int                   `json:"id"`
	Name                   string                `json:"name"`
	Url                    string                `json:"url"`
	Description            string                `json:"description"`
	Notes                  string                `json:"notes"`
	TradeAssocId           string                `json:"tradeAssocId"`
	IsExternal             bool                  `json:"isExternal"`
	CommunityType          string                `json:"CommunityType"`
	ChannelId              string                `json:"ChannelId"`
	OnBoardingInstructions string                `json:"onBoardingInstructions"`
	Created                string                `json:"created"`
	CreatedBy              string                `json:"createdBy"`
	Modified               string                `json:"modified"`
	ModifiedBy             string                `json:"modifiedBy"`
	Sponsors               []SponsorDto          `json:"sponsors"`
	Tags                   []string              `json:"tags"`
	CommunitiesExternal    []RelatedCommunityDto `json:"communitiesExternal"`
	CommunitiesInternal    []RelatedCommunityDto `json:"communitiesInternal"`
}

type SponsorDto struct {
	DisplayName string `json:"displayName"`
	Mail        string `json:"mail"`
}

type RelatedCommunityDto struct {
	ParentCommunityId  int `json:"ParentCommunityId"`
	RelatedCommunityId int `json:"RelatedCommunityId"`
}

type CommunitySponsorsDto struct {
	Id                string `json:"id"`
	CommunityId       string `json:"communityId"`
	UserPrincipalName string `json:"userprincipalname"`
	Created           string `json:"created"`
	CreatedBy         string `json:"createdBy"`
	Modified          string `json:"modified"`
	ModifiedBy        string `json:"modifiedBy"`
}

type CommunityApprovalSystemPostResponseDto struct {
	ItemId string `json:"itemId"`
}

type CommunityApprovalSystemPost struct {
	ApplicationId       string
	ApplicationModuleId string
	Emails              []string
	Subject             string
	Body                string
	RequesterEmail      string
}

type CommunityApprovers struct {
	Id                        int    `json:"id"`
	ApproverUserPrincipalName string `json:"name"`
	Disabled                  bool   `json:"disabled"`
	Created                   string `json:"created"`
	CreatedBy                 string `json:"createdBy"`
	Modified                  string `json:"modified"`
	ModifiedBy                string `json:"modifiedBy"`
}

func AddCommunity(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	var body CommunityDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	param := map[string]interface{}{

		"Name":                   strings.TrimSpace(body.Name),
		"Url":                    body.Url,
		"Description":            body.Description,
		"Notes":                  body.Notes,
		"TradeAssocId":           body.TradeAssocId,
		"CommunityType":          body.CommunityType,
		"ChannelId":              body.ChannelId,
		"OnBoardingInstructions": body.OnBoardingInstructions,
		"CreatedBy":              username,
		"ModifiedBy":             username,
		"Id":                     body.Id,
	}

	result, err := db.CommunitiesInsert(param)
	if err != nil {
		logger.LogException(err)
	}

	id, _ := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
	if err != nil {
		logger.LogException(err)
	}

	for _, s := range body.Sponsors {
		err = db.InsertUser(s.Mail, s.DisplayName, "", "", "")
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sponsorsParam := map[string]interface{}{

			"CommunityId":        id,
			"UserPrincipalName ": s.Mail,
			"CreatedBy":          username,
		}

		_, err := db.CommunitySponsorsInsert(sponsorsParam)

		if err != nil {
			logger.LogException(err)
		}
	}

	deleteParam := map[string]interface{}{

		"ParentCommunityId": id,
	}
	_, err = db.RelatedCommunitiesDelete(deleteParam)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, t := range body.CommunitiesExternal {
		relatedCommunities := map[string]interface{}{
			"ParentCommunityId":   id,
			"RelatedCommunityId ": t.RelatedCommunityId,
		}

		_, err := db.RelatedCommunitiesInsert(relatedCommunities)
		if err != nil {
			logger.LogException(err)
		}
	}

	for _, t := range body.CommunitiesInternal {
		param := map[string]interface{}{
			"ParentCommunityId":   id,
			"RelatedCommunityId ": t.RelatedCommunityId,
		}
		_, err := db.RelatedCommunitiesInsert(param)
		if err != nil {
			logger.LogException(err)
		}
	}

	for _, t := range body.Tags {
		tagsParam := map[string]interface{}{
			"CommunityId": id,
			"Tag ":        t,
		}
		_, err := db.CommunityTagsInsert(tagsParam)
		if err != nil {
			logger.LogException(err)
		}
	}
	if body.Id == 0 {
		requestCommunityApproval(int64(id), logger)
	}

	go getTeamsChannelMembers(body.ChannelId, id)
}

func GetRequestStatusByCommunityId(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id, err := strconv.ParseInt(req["id"], 10, 64)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get project list
	approvalRequests, err := db.CommunityApprovalsSelectById(id)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvalRequests)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetCommunitiesIsExternal(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	isExternal := req["isexternal"]
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	param := map[string]interface{}{

		"isexternal":        isExternal,
		"UserPrincipalName": username,
	}

	communities, err := db.CommunitiesIsexternal(param)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(communities)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func UploadCommunityMembers(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	fmt.Println("Process Community Members List By Excel triggered.")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	counter := 0

	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("fileupload")
	if err != nil {
		logger.LogException(err)
		fmt.Println("Error Retrieving the File")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// check if the file is an excel file
	fileName := strings.Split(handler.Filename, ".")
	if !strings.EqualFold(fileName[len(fileName)-1], "xls") && !strings.EqualFold(fileName[len(fileName)-1], "xlsx") {
		http.Error(w, "The file is not valid.", http.StatusBadRequest)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	xl, _ := xlsxreader.NewReader(fileBytes)
	for row := range xl.ReadRows(xl.Sheets[0]) {
		for _, cell := range row.Cells {
			_, err := mail.ParseAddress(cell.Value)
			if err == nil {
				err = db.Communities_AddMember(id, cell.Value)
				if err == nil {
					counter++
				}
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		NewMembers int `json:"newMembers"`
	}{NewMembers: counter}
	jsonResp, err := json.Marshal(response)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetCommunityOnBoardingInfo(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 0, 64)

	// Get email address of the user
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	switch r.Method {
	case "GET":
		related, _ := db.Communities_Related(id)
		sponsors, _ := db.Community_Sponsors(id)
		info, _ := db.Community_Info(id)
		info.Sponsors = sponsors
		info.Communities = related

		isMember, _ := db.Community_Membership_IsMember(id, username.(string))

		data := map[string]interface{}{
			"IsMember":  isMember,
			"Community": info,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		jsonResp, err := json.Marshal(data)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(jsonResp)

	case "POST":
		db.Community_Onboarding_AddMember(id, username.(string))
	case "DELETE":
		db.Community_Onboarding_RemoveMember(id, username.(string))
	}
}

func GetCommunities(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	result := db.GetCommunities()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResp, err := json.Marshal(result)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(jsonResp)
}

func GetCommunityMembersByCommunityId(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 0, 64)

	result := db.GetCommunityMembers(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResp, err := json.Marshal(result)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(jsonResp)
}

func GetCommunitySponsorsByCommunityId(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id := req["id"]

	// Get project list
	param := map[string]interface{}{

		"CommunityId": id,
	}

	communitySponsors, err := db.CommunitySponsorsSelectByCommunityId(param)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(communitySponsors)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetCommunityTagsByCommunityId(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id := req["id"]

	// Get project list
	param := map[string]interface{}{

		"CommunityId": id,
	}

	communityTags, err := db.CommunityTagsSelectByCommunityId(param)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(communityTags)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func RelatedCommunitiesInsert(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	var body RelatedCommunityDto

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	param := map[string]interface{}{

		"ParentCommunityId":  body.ParentCommunityId,
		"RelatedCommunityId": body.RelatedCommunityId,
	}

	approvers, err := db.RelatedCommunitiesInsert(param)
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

func RelatedCommunitiesDelete(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	var body RelatedCommunityDto

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	param := map[string]interface{}{

		"ParentCommunityId":  body.ParentCommunityId,
		"RelatedCommunityId": body.RelatedCommunityId,
	}

	approvers, err := db.RelatedCommunitiesDelete(param)
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

func GetRelatedCommunitiesByCommunityId(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id := req["id"]

	param := map[string]interface{}{

		"ParentCommunityId": id,
	}

	approvers, err := db.RelatedCommunitiesSelect(param)
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

func GetMyCommunities(w http.ResponseWriter, r *http.Request) {
	// Get email address of the user
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"].(string)

	communities, err := db.MyCommunitites(username)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(communities)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetIManageCommunities(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	params := make(map[string]interface{})
	params["UserPrincipalName"] = username
	communities, err := db.CommunityIManageExecuteSelect(params)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(communities)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetCommunityById(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	communities, err := db.CommunitiesSelectByID(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(communities)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetCommunityApproversById(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	approvers, err := db.GetCommunityApproversById(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvers)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetCommunityApproversList(w http.ResponseWriter, r *http.Request) {
	approvers, err := db.GetCommunityApprovers("community")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvers)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetAllActiveCommunityApprovers(w http.ResponseWriter, r *http.Request) {
	approvers, err := db.GetActiveCommunityApprovers("community")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvers)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func SubmitCommunityApprover(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"].(string)
	var body CommunityApprovers

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.UpdateCommunityApproversById(body.Id, body.Disabled, body.ApproverUserPrincipalName, username, "community")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func requestCommunityApproval(id int64, logger *appinsights_wrapper.TelemetryClient) error {
	communityApprovals := db.PopulateCommunityApproval(id)

	for _, v := range communityApprovals {
		// Create Approval Request
		err := ApprovalSystemRequestCommunity(v, logger)
		if err != nil {
			logger.LogTrace("ID:"+strconv.FormatInt(v.Id, 10)+" "+err.Error(), contracts.Error)
			return err
		}
	}
	return nil
}

func ReprocessRequestCommunityApproval() {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	projectApprovals := db.GetFailedCommunityApprovalRequests()

	for _, v := range projectApprovals {
		err := ApprovalSystemRequestCommunity(v, logger)
		if err != nil {
			logger.LogTrace("ID:"+strconv.FormatInt(v.Id, 10)+" "+err.Error(), contracts.Error)
		}
	}
}

func ApprovalSystemRequestCommunity(data db.CommunityApproval, logger *appinsights_wrapper.TelemetryClient) error {

	url := os.Getenv("APPROVAL_SYSTEM_APP_URL")
	if url != "" {
		url = url + "/api/request"
		ch := make(chan *http.Response)
		// var res *http.Response

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
						<td class="center-table">
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
		)
		body := replacer.Replace(bodyTemplate)

		postParams := CommunityApprovalSystemPost{
			ApplicationId:       os.Getenv("APPROVAL_SYSTEM_APP_ID"),
			ApplicationModuleId: os.Getenv("APPROVAL_SYSTEM_APP_MODULE_COMMUNITY"),
			Emails:              []string{data.ApproverUserPrincipalName},
			Subject:             fmt.Sprintf("[GH-Management] New Community For Approval - %v", data.CommunityName),
			Body:                body,
			RequesterEmail:      data.RequesterUserPrincipalName,
		}

		go getHttpPostResponseStatus(url, postParams, ch, logger)
		r := <-ch
		if r != nil {
			var res CommunityApprovalSystemPostResponseDto
			err := json.NewDecoder(r.Body).Decode(&res)
			if err != nil {
				return err
			}

			db.CommunityApprovalUpdateGUID(data.Id, res.ItemId)
		}
	}
	return nil
}

func getTeamsChannelMembers(channelId string, id int) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	logger.TrackTrace("Community Id: "+fmt.Sprint(id), contracts.Information)
	teamMembers, err := msgraph.GetTeamsMembers(channelId, "")
	if err != nil {
		logger.LogException(err)
	}

	if len(teamMembers) > 0 {
		for _, teamMember := range teamMembers {
			db.Communities_AddMember(id, teamMember.Email)
		}
	}
}

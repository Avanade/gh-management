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

	"main/models"
	db "main/pkg/ghmgmtdb"
	"main/pkg/msgraph"
	"main/pkg/session"

	"github.com/gorilla/mux"
	"github.com/thedatashed/xlsxreader"
)

func CommunityAPIHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	var body models.TypCommunity
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "POST":
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
			log.Println(err.Error())
		}

		id, _ := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
		if err != nil {
			log.Println(err.Error())
		}

		for _, s := range body.Sponsors {
			err = db.InsertUser(s.Mail, s.DisplayName, "", "", "")
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			sponsorsparam := map[string]interface{}{

				"CommunityId":        id,
				"UserPrincipalName ": s.Mail,
				"CreatedBy":          username,
			}

			_, err := db.CommunitySponsorsInsert(sponsorsparam)

			if err != nil {
				log.Println(err.Error())
			}
		}

		deleteparam := map[string]interface{}{

			"ParentCommunityId": id,
		}
		_, err = db.RelatedCommunitiesDelete(deleteparam)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, t := range body.CommunitiesExternal {
			RelatedCommunities := map[string]interface{}{
				"ParentCommunityId":   id,
				"RelatedCommunityId ": t.RelatedCommunityId,
			}

			_, err := db.RelatedCommunitiesInsert(RelatedCommunities)
			if err != nil {
				log.Println(err.Error())
			}
		}

		for _, t := range body.CommunitiesInternal {
			param := map[string]interface{}{
				"ParentCommunityId":   id,
				"RelatedCommunityId ": t.RelatedCommunityId,
			}
			_, err := db.RelatedCommunitiesInsert(param)
			if err != nil {
				log.Println(err.Error())
			}
		}

		for _, t := range body.Tags {
			Tagsparam := map[string]interface{}{
				"CommunityId": id,
				"Tag ":        t,
			}
			_, err := db.CommunityTagsInsert(Tagsparam)
			if err != nil {
				log.Println(err.Error())
			}
		}
		if body.Id == 0 {
			go RequestCommunityApproval(int64(id))
		}

		go func(channelId string) {
			TeamMembers, err := msgraph.GetTeamsMembers(body.ChannelId, "")
			if err != nil {
				log.Println(err.Error())
			}

			if len(TeamMembers) > 0 {
				for _, TeamMember := range TeamMembers {
					db.Communities_AddMember(id, TeamMember.Email)
				}
			}
		}(body.ChannelId)

	case "GET":
		_, err := db.CommunitiesSelectByID(strconv.Itoa(body.Id))
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "PUT":
		param := map[string]interface{}{
			"Id":           body.Id,
			"Name":         body.Name,
			"Url":          body.Url,
			"Description":  body.Description,
			"Notes":        body.Notes,
			"TradeAssocId": body.TradeAssocId,
			"CreatedBy":    body.CreatedBy,
			"ModifiedBy":   body.ModifiedBy,
		}

		_, err := db.CommunitiesUpdate(param)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func MyCommunityAPIHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	var body models.TypCommunity
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
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
			log.Println(err.Error())
		}

		id, _ := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, s := range body.Sponsors {
			err := db.InsertUser(s.Mail, s.DisplayName, "", "", "")
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			sponsorsparam := map[string]interface{}{
				"CommunityId":        id,
				"UserPrincipalName ": s.Mail,
				"CreatedBy":          username,
			}
			_, err = db.CommunitySponsorsInsert(sponsorsparam)
			if err != nil {
				log.Println(err.Error())
			}

		}

		for _, t := range body.Tags {

			Tagsparam := map[string]interface{}{
				"CommunityId": id,
				"Tag ":        t,
			}
			_, err := db.CommunityTagsInsert(Tagsparam)
			if err != nil {
				log.Println(err.Error())
			}

		}

		for _, t := range body.CommunitiesExternal {
			RelatedCommunities := map[string]interface{}{

				"ParentCommunityId":   id,
				"RelatedCommunityId ": t.RelatedCommunityId,
			}
			_, err := db.RelatedCommunitiesInsert(RelatedCommunities)
			if err != nil {
				log.Println(err.Error())
			}
		}

		for _, t := range body.CommunitiesInternal {

			param := map[string]interface{}{

				"ParentCommunityId":   id,
				"RelatedCommunityId ": t.RelatedCommunityId,
			}
			_, err := db.RelatedCommunitiesInsert(param)
			if err != nil {
				log.Println(err.Error())
			}
		}

		if body.Id == 0 {
			go RequestCommunityApproval(int64(id))
		}
	case "GET":
		_, err := db.CommunitiesSelectByID(strconv.Itoa(body.Id))
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "PUT":
		param := map[string]interface{}{

			"Id":           body.Id,
			"Name":         body.Name,
			"Url":          body.Url,
			"Description":  body.Description,
			"Notes":        body.Notes,
			"TradeAssocId": body.TradeAssocId,
			"CreatedBy":    body.CreatedBy,
			"ModifiedBy":   body.ModifiedBy,
		}
		_, err := db.CommunitiesUpdate(param)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetRequestStatusByCommunity(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	// Get project list
	params := make(map[string]interface{})
	params["Id"] = id
	projects, err := db.CommunityApprovalsSelectById(params)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(projects)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetCommunitiesIsexternal(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	isexternal := req["isexternal"]
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	param := map[string]interface{}{

		"isexternal":        isexternal,
		"UserPrincipalName": username,
	}

	Communities, err := db.CommunitiesIsexternal(param)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(Communities)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func CommunityInitCommunityType(w http.ResponseWriter, r *http.Request) {
	_, err := db.CommunitiesInitCommunityType(nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ProcessCommunityMembersListExcel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Process Community Members List By Excel triggered.")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	counter := 0

	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("fileupload")
	if err != nil {
		log.Println(err.Error())
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
		log.Println(err.Error())
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
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func RequestCommunityApproval(id int64) error {
	communityApprovals := db.PopulateCommunityApproval(id)

	for _, v := range communityApprovals {
		err := ApprovalSystemRequestCommunity(v)
		if err != nil {
			log.Println("ID:" + strconv.FormatInt(v.Id, 10) + " " + err.Error())
			return err
		}
	}
	return nil
}

func ReprocessRequestCommunityApproval() {
	projectApprovals := db.GetFailedCommunityApprovalRequests()

	for _, v := range projectApprovals {
		go ApprovalSystemRequestCommunity(v)
	}
}

func ApprovalSystemRequestCommunity(data models.TypCommunityApprovals) error {

	url := os.Getenv("APPROVAL_SYSTEM_APP_URL")
	if url != "" {
		url = url + "/request"
		ch := make(chan *http.Response)
		// var res *http.Response

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

		postParams := models.TypApprovalSystemPost{
			ApplicationId:       os.Getenv("APPROVAL_SYSTEM_APP_ID"),
			ApplicationModuleId: os.Getenv("APPROVAL_SYSTEM_APP_MODULE_COMMUNITY"),
			Email:               data.ApproverUserPrincipalName,
			Subject:             fmt.Sprintf("[GH-Management] New Community For Approval - %v", data.CommunityName),
			Body:                body,
			RequesterEmail:      data.RequesterUserPrincipalName,
		}

		go getHttpPostResponseStatus(url, postParams, ch)
		r := <-ch
		if r != nil {
			var res models.TypApprovalSystemPostResponse
			err := json.NewDecoder(r.Body).Decode(&res)
			if err != nil {
				log.Println(err.Error())
				return err
			}

			db.CommunityApprovalUpdateGUID(data.Id, res.ItemId)
		}
	}
	return nil
}

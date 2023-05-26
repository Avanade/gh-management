package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	models "main/models"
	ghmgmt "main/pkg/ghmgmtdb"
	"main/pkg/msgraph"
	session "main/pkg/session"
	"main/pkg/sql"
	comm "main/routes/pages/community"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		fmt.Println(body)
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

		result, err := ghmgmt.CommunitiesInsert(param)
		if err != nil {
			fmt.Println(err)
		}
		id, _ := strconv.Atoi(fmt.Sprint(result[0]["Id"]))

		if err != nil {
			fmt.Println(err)
		}

		for _, s := range body.Sponsors {
			errIU := ghmgmt.InsertUser(s.Mail, s.DisplayName, "", "", "")
			if errIU != nil {
				http.Error(w, errIU.Error(), http.StatusInternalServerError)
				return
			}

			sponsorsparam := map[string]interface{}{

				"CommunityId":        id,
				"UserPrincipalName ": s.Mail,
				"CreatedBy":          username,
			}

			_, err := ghmgmt.CommunitySponsorsInsert(sponsorsparam)

			if err != nil {
				fmt.Println(err)

			}

		}

		deleteparam := map[string]interface{}{

			"ParentCommunityId": id,
		}
		_, error := ghmgmt.RelatedCommunitiesDelete(deleteparam)
		if err != nil {

			fmt.Println(error)
		}

		for _, t := range body.CommunitiesExternal {

			RelatedCommunities := map[string]interface{}{

				"ParentCommunityId":   id,
				"RelatedCommunityId ": t.RelatedCommunityId,
			}

			_, err := ghmgmt.RelatedCommunitiesInsert(RelatedCommunities)
			if err != nil {

				fmt.Println(err)
			}

		}

		for _, t := range body.CommunitiesInternal {

			param := map[string]interface{}{

				"ParentCommunityId":   id,
				"RelatedCommunityId ": t.RelatedCommunityId,
			}
			_, err := ghmgmt.RelatedCommunitiesInsert(param)
			if err != nil {

				fmt.Println(err)
			}

		}
		for _, t := range body.Tags {

			Tagsparam := map[string]interface{}{

				"CommunityId": id,
				"Tag ":        t,
			}
			_, err := ghmgmt.CommunityTagsInsert(Tagsparam)
			if err != nil {

				fmt.Println(err)
			}

		}
		if body.Id == 0 {
			go comm.RequestCommunityApproval(int64(id))
		}

		go func(channelId string) {
			TeamMembers, err := msgraph.GetTeamsMembers(body.ChannelId, "")
			if err != nil {

				fmt.Println(err)
				return
			}
			if len(TeamMembers) > 0 {

				for _, TeamMember := range TeamMembers {

					ghmgmt.Communities_AddMember(id, TeamMember.Email)
				}
			}

		}(body.ChannelId)

	case "GET":
		param := map[string]interface{}{

			"Id": body.Id,
		}
		_, err := ghmgmt.CommunitiesSelectByID(param)
		if err != nil {
			fmt.Println(err)
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

		_, err := ghmgmt.CommunitiesUpdate(param)
		if err != nil {
			fmt.Println(err)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		fmt.Println(body)
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

		result, err := ghmgmt.CommunitiesInsert(param)
		if err != nil {
			fmt.Println(err)
		}
		id, _ := strconv.Atoi(fmt.Sprint(result[0]["Id"]))

		if err != nil {
			fmt.Println(err)
		}

		for _, s := range body.Sponsors {
			errIU := ghmgmt.InsertUser(s.Mail, s.DisplayName, "", "", "")
			if errIU != nil {
				http.Error(w, errIU.Error(), http.StatusInternalServerError)
				return
			}

			sponsorsparam := map[string]interface{}{

				"CommunityId":        id,
				"UserPrincipalName ": s.Mail,
				"CreatedBy":          username,
			}
			_, err := ghmgmt.CommunitySponsorsInsert(sponsorsparam)
			if err != nil {
				fmt.Println(err)

			}

		}

		for _, t := range body.Tags {

			Tagsparam := map[string]interface{}{

				"CommunityId": id,
				"Tag ":        t,
			}
			_, err := ghmgmt.CommunityTagsInsert(Tagsparam)
			if err != nil {

				fmt.Println(err)
			}

		}

		for _, t := range body.CommunitiesExternal {

			RelatedCommunities := map[string]interface{}{

				"ParentCommunityId":   id,
				"RelatedCommunityId ": t.RelatedCommunityId,
			}
			_, err := ghmgmt.RelatedCommunitiesInsert(RelatedCommunities)
			if err != nil {

				fmt.Println(err)
			}

		}

		for _, t := range body.CommunitiesInternal {

			param := map[string]interface{}{

				"ParentCommunityId":   id,
				"RelatedCommunityId ": t.RelatedCommunityId,
			}
			_, err := ghmgmt.RelatedCommunitiesInsert(param)
			if err != nil {

				fmt.Println(err)
			}

		}

		if body.Id == 0 {
			go comm.RequestCommunityApproval(int64(id))
		}
	case "GET":
		param := map[string]interface{}{

			"Id": body.Id,
		}
		_, err := ghmgmt.CommunitiesSelectByID(param)
		if err != nil {
			fmt.Println(err)
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
		_, err := ghmgmt.CommunitiesUpdate(param)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func GetRequestStatusByCommunity(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

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

	// Get project list
	params := make(map[string]interface{})
	params["Id"] = id
	projects, err := ghmgmt.CommunityApprovalsSelectById(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(projects)
	if err != nil {
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

		"isexternal":        isexternal,
		"UserPrincipalName": username,
	}

	Communities, err := ghmgmt.CommunitiesIsexternal(param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(Communities)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}
func CommunityInitCommunityType(w http.ResponseWriter, r *http.Request) {
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = ghmgmt.CommunitiesInitCommunityType(nil)
	if err != nil {
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

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	xl, _ := xlsxreader.NewReader(fileBytes)
	for row := range xl.ReadRows(xl.Sheets[0]) {
		for _, cell := range row.Cells {
			_, err := mail.ParseAddress(cell.Value)
			if err == nil {
				err = ghmgmt.Communities_AddMember(id, cell.Value)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

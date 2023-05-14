package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	ghmgmt "main/pkg/ghmgmtdb"
	"main/pkg/msgraph"
	session "main/pkg/session"
	"main/pkg/sql"
	comm "main/routes/pages/community"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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
	cp := sql.ConnectionParam{

		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)
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

		result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Communities_Insert", param)
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
			_, err := db.ExecuteStoredProcedure("dbo.PR_CommunitySponsors_Insert", sponsorsparam)
			if err != nil {
				fmt.Println(err)

			}

		}

		deleteparam := map[string]interface{}{

			"ParentCommunityId": id,
		}
		_, error := db.ExecuteStoredProcedure("PR_RelatedCommunities_Delete", deleteparam)
		if err != nil {

			fmt.Println(error)
		}

		for _, t := range body.CommunitiesExternal {

			RelatedCommunities := map[string]interface{}{

				"ParentCommunityId":   id,
				"RelatedCommunityId ": t.RelatedCommunityId,
			}
			_, err := db.ExecuteStoredProcedure("PR_RelatedCommunities_Insert", RelatedCommunities)
			if err != nil {

				fmt.Println(err)
			}

		}

		for _, t := range body.CommunitiesInternal {

			param := map[string]interface{}{

				"ParentCommunityId":   id,
				"RelatedCommunityId ": t.RelatedCommunityId,
			}
			_, err := db.ExecuteStoredProcedure("PR_RelatedCommunities_Insert", param)
			if err != nil {

				fmt.Println(err)
			}

		}
		for _, t := range body.Tags {

			Tagsparam := map[string]interface{}{

				"CommunityId": id,
				"Tag ":        t,
			}
			_, err := db.ExecuteStoredProcedure("PR_CommunityTags_Insert", Tagsparam)
			if err != nil {

				fmt.Println(err)
			}

		}
		if body.Id == 0 {
			go comm.RequestCommunityApproval(int64(id))
		}

		go func(channelId string) {
			TeamMembers, _ := msgraph.GetTeamsMembers(body.ChannelId, "")
			if len(TeamMembers) > 0 {

				for _, TeamMember := range TeamMembers {
					fmt.Println(TeamMember.Email)

					ghmgmt.Communities_AddMember(id, TeamMember.Email)
				}
			}

		}(body.ChannelId)

	case "GET":
		param := map[string]interface{}{

			"Id": body.Id,
		}
		_, err := db.ExecuteStoredProcedure("dbo.PR_Communities_select_byID", param)
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
		_, err := db.ExecuteStoredProcedure("dbo.PR_Communities_Update", param)
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
	cp := sql.ConnectionParam{

		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)
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

		result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Communities_Insert", param)
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
			_, err := db.ExecuteStoredProcedure("dbo.PR_CommunitySponsors_Insert", sponsorsparam)
			if err != nil {
				fmt.Println(err)

			}

		}

		for _, t := range body.Tags {

			Tagsparam := map[string]interface{}{

				"CommunityId": id,
				"Tag ":        t,
			}
			_, err := db.ExecuteStoredProcedure("PR_CommunityTags_Insert", Tagsparam)
			if err != nil {

				fmt.Println(err)
			}

		}

		for _, t := range body.CommunitiesExternal {

			RelatedCommunities := map[string]interface{}{

				"ParentCommunityId":   id,
				"RelatedCommunityId ": t.RelatedCommunityId,
			}
			_, err := db.ExecuteStoredProcedure("PR_RelatedCommunities_Insert", RelatedCommunities)
			if err != nil {

				fmt.Println(err)
			}

		}

		for _, t := range body.CommunitiesInternal {

			param := map[string]interface{}{

				"ParentCommunityId":   id,
				"RelatedCommunityId ": t.RelatedCommunityId,
			}
			_, err := db.ExecuteStoredProcedure("PR_RelatedCommunities_Insert", param)
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
		_, err := db.ExecuteStoredProcedure("dbo.PR_Communities_select_byID", param)
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
		_, err := db.ExecuteStoredProcedure("dbo.PR_Communities_Update", param)
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
	projects, err := db.ExecuteStoredProcedureWithResult("PR_CommunityApprovals_Select_ById", params)
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

	Communities, err := db.ExecuteStoredProcedureWithResult("PR_Communities_Isexternal", param)
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

	_, err = db.ExecuteStoredProcedure("dbo.PR_Communities_InitCommunityType", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
}
func ConnectDb() *sql.DB {
	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	return db
}

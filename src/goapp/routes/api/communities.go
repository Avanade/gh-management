package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	ghmgmt "main/pkg/ghmgmtdb"
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

			"Name":         strings.TrimSpace(body.Name),
			"Url":          body.Url,
			"Description":  body.Description,
			"Notes":        body.Notes,
			"TradeAssocId": body.TradeAssocId,
			"IsExternal":   body.IsExternal,
			"CreatedBy":    username,
			"ModifiedBy":   username,
			"Id":           body.Id,
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

func ConnectDb() *sql.DB {
	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	return db
}

package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	session "main/pkg/session"
	"main/pkg/sql"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func CommunitySponsorsAPIHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	var body models.TypCommunitySponsors
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cp := sql.ConnectionParam{

		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)
	switch r.Method {
	case "POST":
		param := map[string]interface{}{

			"CommunityId":       body.CommunityId,
			"UserPrincipalName": body.UserPrincipalName,
			"CreatedBy":         username,
		}
		_, err := db.ExecuteStoredProcedure("dbo.PR_CommunitySponsors_Insert", param)
		if err != nil {
			fmt.Println(err)
		}
	case "GET":
		//	param := map[string]interface{}{}
		_, err := db.ExecuteStoredProcedure("dbo.PR_CommunitySponsors_Select", nil)
		if err != nil {
			fmt.Println(err)
		}

	case "PUT":
		param := map[string]interface{}{
			"CommunityId":       body.CommunityId,
			"UserPrincipalName": body.UserPrincipalName,
			"CreatedBy":         username,
		}
		_, err := db.ExecuteStoredProcedure("dbo.PR_CommunitySponsors_Update", param)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func CommunitySponsorsPerCommunityId(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]
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
	param := map[string]interface{}{

		"CommunityId": id,
	}
	CommunitySponsors, err := db.ExecuteStoredProcedureWithResult("PR_CommunitySponsors_Select_By_CommunityId", param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(CommunitySponsors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Printf(projects)
	w.Write(jsonResp)
}

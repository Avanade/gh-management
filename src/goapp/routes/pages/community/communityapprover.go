package routes

import (
	"encoding/json"

	session "main/pkg/session"
	"main/pkg/sql"
	template "main/pkg/template"
	"net/http"
	"os"

	"fmt"
	models "main/models"

	"github.com/gorilla/mux"
)

func CommunityApproverHandler(w http.ResponseWriter, r *http.Request) {

	template.UseTemplate(&w, r, "/community/communityapprovers", nil)
}

func GetCommunityApproversById(w http.ResponseWriter, r *http.Request) {
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

	param := map[string]interface{}{

		"Id": id,
	}

	approvers, err := db.ExecuteStoredProcedureWithResult("PR_CommunityApproversList_select_byId", param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetCommunityApproversList(w http.ResponseWriter, r *http.Request) {

	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	approvers, err := db.ExecuteStoredProcedureWithResult("PR_CommunityApproversList_select", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetAllActiveCommunityApprovers(w http.ResponseWriter, r *http.Request) {

	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	approvers, err := db.ExecuteStoredProcedureWithResult("PR_CommunityApproversList_SelectAllActive", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func CommunityApproversListUpdate(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	var body models.TypCommunityApprovers

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)

		return
	}

	cp := sql.ConnectionParam{

		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}
	db, _ := sql.Init(cp)
	defer db.Close()

	param1 := map[string]interface{}{

		"ApproverUserPrincipalName": body.ApproverUserPrincipalName,
		"Disabled":                  body.Disabled,
		"CreatedBy":                 username,
		"ModifiedBy":                username,
		"Id":                        body.Id,
	}

	__, err := db.ExecuteStoredProcedureWithResult("PR_CommunityApproversList_Insert", param1)
	if err != nil {
		fmt.Println(__)
		fmt.Println(err)

	}

}

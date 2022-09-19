package routes

import (
	"encoding/json"
	"fmt"

	//session "main/pkg/session"
	"main/pkg/sql"
	template "main/pkg/template"
	"net/http"
	"os"

	//models "main/models"
	//"fmt"
	"github.com/gorilla/mux"
)

func CommunitylistHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CommunitylistHandler")
	template.UseTemplate(&w, r, "community/communitylist", nil)
}

// func CommunityApproverHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("CommunityApproverHandler")
// 	template.UseTemplate(&w, r, "community/communityapprovers", nil)
// }

func GetUserCommunitylist(w http.ResponseWriter, r *http.Request) {

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

	Communities, err := db.ExecuteStoredProcedureWithResult("PR_Communities_select", nil)
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

func GetUserCommunity(w http.ResponseWriter, r *http.Request) {
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

	Communities, err := db.ExecuteStoredProcedureWithResult("PR_Communities_select_byID", param)
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

package routes

import (
	"encoding/json"
	//session "main/pkg/session"
	"main/pkg/sql"
	template "main/pkg/template"
	"net/http"
	"os"
	"fmt"
	//"github.com/gorilla/mux"
)

func CommunitylistHandler(w http.ResponseWriter, r *http.Request) {
 fmt.Println("test")
	template.UseTemplate(&w, r, "community/communitylist", nil)
}

func GetUserCommunitylist(w http.ResponseWriter, r *http.Request) {
	// Get email address of the user
//	sessionaz, _ := session.Store.Get(r, "auth-session")
//	iprofile := sessionaz.Values["profile"]
//	profile := iprofile.(map[string]interface{})
//	username := profile["preferred_username"]
fmt.Println("test2")
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
 
	projects, err := db.ExecuteStoredProcedureWithResult("PR_Communities_select", nil)
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
	//fmt.Printf(projects)
	w.Write(jsonResp)
}
package routes

import (
	"encoding/json"
	"main/pkg/sql"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	template "main/pkg/template"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	//users := db.GetUsersWithGithub()
	template.UseTemplate(&w, r, "search/search", nil)

}

func GetSearchResults(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	searchText := req["searchText"]

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

	// Get Searh List from SP
	params := make(map[string]interface{})
	params["searchText"] = searchText
	searchResults, err := db.ExecuteStoredProcedureWithResult("PR_Search_communities_projects_users", params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(searchResults)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetAllResults(w http.ResponseWriter, r *http.Request) {
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

	// Get Searh List from SP
	searchResults, err := db.ExecuteStoredProcedureWithResult("PR_SearchAll_communities_projects_users", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(searchResults)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

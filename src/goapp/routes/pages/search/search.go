package routes

import (
	"encoding/json"
	session "main/pkg/session"
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

	searchText := r.URL.Query().Get("search")
	offSet := req["offSet"]
	rowCount := req["rowCount"]
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
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

	// Get Search List from SP
	//searchTextParam := make(map[string]interface{})
	params := map[string]interface{}{
		"searchText":    searchText,
		"offSet":        offSet,
		"rowCount":      rowCount,
		"userprincipal": username,
	}

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

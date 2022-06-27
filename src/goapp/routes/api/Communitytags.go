package routes

import (
	"encoding/json"
	"main/pkg/sql"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func CommunityTagPerCommunityId(w http.ResponseWriter, r *http.Request) {
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
	CommunityTags, err := db.ExecuteStoredProcedureWithResult("PR_CommunityTags_Select_By_CommunityId", param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(CommunityTags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Printf(projects)
	w.Write(jsonResp)
}

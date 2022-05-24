package route

import (
	"encoding/json"
	"fmt"
	"main/pkg/session"
	"main/pkg/sql"
	template "main/pkg/template"
	"net/http"
	"os"
)

func MyRequestsHandler(w http.ResponseWriter, r *http.Request) {

	// Connect to database
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("APPROVALSYSTEMDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get email of the user
	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	u := session.Values["profile"]
	profile, ok := u.(map[string]interface{})
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userEmail := fmt.Sprintf("%s", profile["preferred_username"])

	// Get approval requests assigned to the user
	params := make(map[string]interface{})
	params["CreatedBy"] = userEmail
	items, err := db.ExecuteStoredProcedureWithResult("PR_Items_Select_ByCreatedBy", params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to struct
	var homeData TypHomeData
	for _, v := range items {
		if v["IsApproved"] == nil {
			homeData.Pending = append(homeData.Pending, itemMapper(v, false))
		} else {
			homeData.Approved = append(homeData.Approved, itemMapper(v, true))
		}
	}

	b, err := json.Marshal(homeData)
	if err != nil {
		fmt.Println(err)
		return
	}

	template.UseTemplate(&w, r, "myrequests", string(b))
}

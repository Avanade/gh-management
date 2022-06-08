package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	session "main/pkg/session"
	"main/pkg/sql"
	"net/http"
	"os"
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

// func ConnectDb() *sql.DB {
// 	cp := sql.ConnectionParam{
// 		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
// 	}

// 	db, _ := sql.Init(cp)

// 	return db
// }

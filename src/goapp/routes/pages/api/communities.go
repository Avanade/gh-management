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

func CommunityAPIHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	var body models.TypCommunity
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

			"Name":         body.Name,
			"Url":          body.Url,
			"Description":  body.Description,
			"Notes":        body.Notes,
			"TradeAssocId": body.TradeAssocId,
			"CreatedBy":    username,
			"ModifiedBy":   username,
		}

		_, err := db.ExecuteStoredProcedure("dbo.PR_Communities_Insert", param)
		if err != nil {
			fmt.Println(err)
		}

		for _, s := range body.Sponsors {

			sponsorsparam := map[string]interface{}{

				"CommunityId":        1,
				"UserPrincipalName ": s,
				"CreatedBy":          username,
			}
			_, err := db.ExecuteStoredProcedure("dbo.PR_CommunitySponsors_Insert", sponsorsparam)
			if err != nil {
				fmt.Println(err)
			}
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

func ConnectDb() *sql.DB {
	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	return db
}

package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	ghmgmt "main/pkg/ghmgmtdb"
	session "main/pkg/session"
	"main/pkg/sql"
	"net/http"
	"os"
	"strconv"
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
		fmt.Println(err)
		fmt.Println(body)
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
			"IsExternal":   body.IsExternal,
			"CreatedBy":    username,
			"ModifiedBy":   username,
			"Id":           body.Id,
		}

		result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Communities_Insert", param)
		if err != nil {
			fmt.Println(err)
		}
		id, _ := strconv.Atoi(fmt.Sprint(result[0]["Id"]))

		if err != nil {
			fmt.Println(err)
		}

		for _, s := range body.Sponsors {
			errIU := ghmgmt.InsertUser(s.Mail, s.DisplayName, "", "", "")
			if errIU != nil {
				http.Error(w, errIU.Error(), http.StatusInternalServerError)
				return
			}
			sponsorsparam := map[string]interface{}{

				"CommunityId":        id,
				"UserPrincipalName ": s.DisplayName,
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

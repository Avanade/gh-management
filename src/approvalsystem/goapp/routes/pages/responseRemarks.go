package routes

import (
	"fmt"
	session "main/pkg/session"
	"main/pkg/sql"
	template "main/pkg/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func ResponseRemarksHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})

	params := mux.Vars(r)

	appGuid := params["appGuid"]
	appModuleGuid := params["appModuleGuid"]
	itemGuid := params["itemGuid"]
	isApproved := params["isApproved"]

	username := profile["preferred_username"]

	db, err := sql.Init(sql.ConnectionParam{ConnectionString: os.Getenv("APPROVALSYSTEMDB_CONNECTION_STRING")})
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
	}
	sqlParamsIsAuth := map[string]interface{}{
		"ApplicationId":       appGuid,
		"ApplicationModuleId": appModuleGuid,
		"ItemId":              itemGuid,
		"ApproverEmail":       username,
	}

	sqlParamsItems := map[string]interface{}{
		"Id": itemGuid,
	}

	resIsAuth, err := db.ExecuteStoredProcedureWithResult("PR_RESPONSE_IsAuthorized", sqlParamsIsAuth)

	isAuth := resIsAuth[0]["IsAuthorized"]
	if isAuth == "0" {
		template.UseTemplate(&w, r, "Unauthorized", nil)
	} else {
		resItems, _ := db.ExecuteStoredProcedureWithResult("PR_Items_Select_ById", sqlParamsItems)

		ApproveUrl := "/processRemarks"
		RejectUrl := "/processRemarks"
		buttons := map[string]string{
			"ApproveUrl": ApproveUrl,
			"RejectUrl":  RejectUrl,
		}
		data := map[string]interface{}{
			"isApproved": isApproved,
			"data":       resItems[0],
			"buttons":    buttons,
		}

		fmt.Println(resItems[0]["Id"])
		template.UseTemplate(&w, r, "responseRemarks", data)

	}
}

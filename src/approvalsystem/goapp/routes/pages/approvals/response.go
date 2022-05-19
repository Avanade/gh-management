package route

import (
	"encoding/json"
	"fmt"
	"main/models"
	session "main/pkg/session"
	"main/pkg/sql"
	template "main/pkg/template"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func ResponseHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
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
		defer db.Close()

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
			isProcessed := resIsAuth[0]["IsApproved"]
			fmt.Println(isProcessed)
			if isProcessed != nil {
				var text string
				if isProcessed == true {
					text = "approved"
				} else {
					text = "rejected"
				}
				data := map[string]interface{}{
					"response": text,
				}
				template.UseTemplate(&w, r, "AlreadyProcessed", data)
			} else {
				resItems, _ := db.ExecuteStoredProcedureWithResult("PR_Items_Select_ById", sqlParamsItems)

				data := map[string]interface{}{
					"ApplicationId":       appGuid,
					"ApplicationModuleId": appModuleGuid,
					"ItemId":              itemGuid,
					"ApproverEmail":       username,
					"isApproved":          isApproved,
					"data":                resItems[0],
				}
				template.UseTemplate(&w, r, "response", data)
			}

		}
	}
}

func ProcessResponseHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// Decode payload
		var req models.TypRequestProcess
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

		// Validate payload
		params := make(map[string]interface{})
		params["ApplicationId"] = req.ApplicationId
		params["ApplicationModuleId"] = req.ApplicationModuleId
		params["ItemId"] = req.ItemId
		params["ApproverEmail"] = req.ApproverEmail
		verification, err := db.ExecuteStoredProcedureWithResult("PR_Items_IsValid", params)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if verification[0]["IsValid"] == "1" {
			for k := range params {
				delete(params, k)
			}
			isApproved, _ := strconv.ParseBool(req.IsApproved)
			params["Id"] = req.ItemId
			params["IsApproved"] = isApproved
			params["ApproverRemarks"] = req.Remarks
			fmt.Println("isApproved", isApproved)
			result, err := db.ExecuteStoredProcedure("PR_Items_Update_Response", params)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println(result)
			w.WriteHeader(200)
			return
		} else {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}
}

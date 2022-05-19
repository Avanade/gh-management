package route

import (
	"encoding/json"
	"fmt"
	"main/models"
	"main/pkg/email"
	"main/pkg/sql"
	"net/http"
	"os"
)

func ApprovalRequestHandler(w http.ResponseWriter, r *http.Request) {

	// Decode payload
	var req models.TypRequestApproval
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
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
	verification, err := db.ExecuteStoredProcedureWithResult("PR_Application_IsValid", params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if verification[0]["IsValid"] == "1" {
		for k := range params {
			delete(params, k)
		}

		// Get application module
		params["Id"] = req.ApplicationModuleId
		appModule, err := db.ExecuteStoredProcedureWithResult("PR_ApplicationModules_Select_ById", params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var baseResponseUrl string
		if appModule[0]["RequireRemarks"].(bool) {
			baseResponseUrl = "responseRemarks"
		} else {
			baseResponseUrl = "response"
		}

		// Add item to database
		for k := range params {
			delete(params, k)
		}
		params["ApplicationModuleId"] = req.ApplicationModuleId
		params["ApproverEmail"] = req.Email
		params["Subject"] = req.Subject
		params["Body"] = req.Body
		item, err := db.ExecuteStoredProcedureWithResult("PR_Items_Insert", params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Compose email message
		emailBodyData := email.TypEmailData{
			Subject:     req.Subject,
			Body:        req.Body,
			ApproveText: fmt.Sprintf("%s", appModule[0]["ApproveText"]),
			RejectText:  fmt.Sprintf("%s", appModule[0]["RejectText"]),
			ApproveUrl:  fmt.Sprintf("%s/%s/%s/%s/%s/1", os.Getenv("HOME_URL"), baseResponseUrl, req.ApplicationId, req.ApplicationModuleId, item[0]["Id"]),
			RejectUrl:   fmt.Sprintf("%s/%s/%s/%s/%s/0", os.Getenv("HOME_URL"), baseResponseUrl, req.ApplicationId, req.ApplicationModuleId, item[0]["Id"]),
		}

		emailBody, err := email.ComposeEmail(emailBodyData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send email to approver
		emailData := email.TypEmailMessage{
			To:      req.Email,
			Subject: req.Subject,
			Body:    emailBody,
		}

		_, err = email.SendEmail(emailData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepare response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["itemId"] = fmt.Sprintf("%v", item[0]["Id"])
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonResp)

	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

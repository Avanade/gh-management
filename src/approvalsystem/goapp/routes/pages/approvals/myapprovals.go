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

func MyApprovalsHandler(w http.ResponseWriter, r *http.Request) {

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
	approverEmail := fmt.Sprintf("%s", profile["preferred_username"])

	// Get approval requests assigned to the user
	params := make(map[string]interface{})
	params["ApproverEmail"] = approverEmail
	items, err := db.ExecuteStoredProcedureWithResult("PR_Items_Select_ByApproverEmail", params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to struct
	var homeData TypHomeData
	for _, v := range items {
		if v["IsApproved"] == nil {
			homeData.Pending = append(homeData.Pending, itemMapper(v, false))
			fmt.Println(v.Created)
		} else {
			homeData.Approved = append(homeData.Approved, itemMapper(v, true))
		}
	}

	b, err := json.Marshal(homeData)
	if err != nil {
		fmt.Println(err)
		return
	}

	template.UseTemplate(&w, r, "myapprovals", string(b))
}

func itemMapper(item map[string]interface{}, isApproved bool) TypItem {
	var approveUrl string
	var rejectUrl string

	if !isApproved {
		approveUrl = fmt.Sprintf("/response/%s/%s/%s/1", item["ApplicationId"], item["ApplicationModuleId"], item["ItemId"])
		rejectUrl = fmt.Sprintf("/response/%s/%s/%s/0", item["ApplicationId"], item["ApplicationModuleId"], item["ItemId"])
	}

	return TypItem{
		Application:     item["Application"],
		ApproverRemarks: item["ApproverRemarks"],
		Body:            item["Body"],
		Created:         item["Created"],
		DateResponded:   item["DateResponded"],
		DateSent:        item["DateSent"],
		IsApproved:      item["IsApproved"],
		Module:          item["Module"],
		Subject:         item["Subject"],
		ApproveText:     item["ApproveText"],
		RejectText:      item["RejectText"],
		ApproveUrl:      approveUrl,
		RejectUrl:       rejectUrl,
		Show:            false,
	}
}

type TypHomeData struct {
	Approved []TypItem
	Pending  []TypItem
}

type TypItem struct {
	Application     interface{}
	ApproverRemarks interface{}
	Body            interface{}
	Created         interface{}
	DateResponded   interface{}
	DateSent        interface{}
	IsApproved      interface{}
	Module          interface{}
	Subject         interface{}
	ApproveText     interface{}
	RejectText      interface{}
	ItemId          interface{}
	ApproveUrl      interface{}
	RejectUrl       interface{}
	Show            bool
}

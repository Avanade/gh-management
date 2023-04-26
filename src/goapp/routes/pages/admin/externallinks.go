package routes

import (
	"encoding/json"

	// db "main/pkg/ghmgmtdb"
	"fmt"
	models "main/models"
	session "main/pkg/session"
	"main/pkg/sql"
	template "main/pkg/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func ExternalLinksHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/externallinks", nil)
}
func ExternalLinksForm(w http.ResponseWriter, r *http.Request) {
	// template.UseTemplate(&w, r, "admin/externallinks/form", nil)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	action := vars["action"]
	template.UseTemplate(&w, r, "admin/externallinks/form", struct {
		Id     int
		Action string
	}{
		Id:     id,
		Action: strings.Title(action),
	})
}

func GetExternalLinks(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])
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

	param := map[string]interface{}{
		"CreatedBy": username,
	}

	ExternalLinks, err := db.ExecuteStoredProcedureWithResult("PR_ExternalLinks_Select", param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(ExternalLinks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetExternalLinksById(w http.ResponseWriter, r *http.Request) {
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

	param := map[string]interface{}{
		"Id": id,
	}

	ExternalLinks, err := db.ExecuteStoredProcedureWithResult("PR_ExternalLinks_SelectById", param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(ExternalLinks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

// func ExternalLinksForm(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, _ := strconv.Atoi(vars["id"])
// 	action := vars["action"]
// 	template.UseTemplate(&w, r, "admin/contributionareas/form", struct {
// 		Id     int
// 		Action string
// 	}{
// 		Id:     id,
// 		Action: strings.Title(action),
// 	})
// }

func GetExternalLinksByCategory(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	Category := req["Category"]

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])
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

	param := map[string]interface{} {
		"Category": Category,
		"CreatedBy": username,
	}

	ExternalLinks, err := db.ExecuteStoredProcedureWithResult("PR_ExternalLinks_SelectByCategory", param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(ExternalLinks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

// Add External Links
func CreateExternalLinks(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	var body models.TypExternalLinks

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}
	db, _ := sql.Init(cp)
	defer db.Close()

	param := map[string]interface{}{
		"SVGName":   body.SVGName,
		"IconSVG":   body.IconSVG,
		"Category":  body.Category,
		"CreatedBy": username,
		"Enabled" :  body.Enabled,
	}

	__, err := db.ExecuteStoredProcedureWithResult("PR_ExternalLinks_Insert", param)
	if err != nil {
		fmt.Println(__)
		fmt.Println(err)

	}
}

func UpdateExternalLinks(w http.ResponseWriter, r *http.Request) {

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	var body models.TypCategory

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)

		return
	}

	cp := sql.ConnectionParam{

		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}
	db, err := sql.Init(cp)

	params := map[string]interface{}{

		"Name":       body.Name,
		"CreatedBy":  username,
		"ModifiedBy": username,
		"Id":         body.Id,
	}

	_, err2 := db.ExecuteStoredProcedureWithResult("dbo.PR_ExternalLinks_Update", params)
	if err != nil {
		fmt.Println(err2)

		return
	}

}

func ExternalLinksDelete(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]



	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(cp)

	param := map[string]interface{}{
		"Id": id,
	}

	_, error := db.ExecuteStoredProcedure("PR_ExternalLinks_Delete", param)
	if err != nil {
		fmt.Println(error)
	}
}

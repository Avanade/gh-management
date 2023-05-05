package routes

import (
	"encoding/json"
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

	var externalLinks []models.TypMenu

	externalLinks = append(externalLinks, models.TypMenu{Name: "Tech Community Calendar", Url: "https://techcommunitycalendar.com/", IconPath: "/public/icons/ExternalLinks/arrow-trending-up.svg", External: true})


	data:= 

	tmpl := template.Must(
		template.ParseFiles("admin/externallinks/form",
			fmt.Sprintf("admin/externallinks/%v.html", page)))
	return tmpl.Execute(*w, data)

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

func GetExternalLinksByCategory(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	category := req["Category"]

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
		"Category":  category,
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

func CreateExternalLinks(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	var data models.TypExternalLinks

	err := json.NewDecoder(r.Body).Decode(&data)
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

	params := map[string]interface{}{
		"SVGName":   data.SVGName,
		"IconSVG":   data.IconSVG,
		"Hyperlink": data.Hyperlink,
		"LinkName":  data.LinkName,
		"Category":  data.Category,
		"Enabled":   data.Enabled,
		"CreatedBy": username,
	}

	__, err := db.ExecuteStoredProcedure("PR_ExternalLinks_Insert", params)
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
	db, err := sql.Init(cp)
	

	params := map[string]interface{}{
		"Id":         body.Id,
		"SVGName":    body.SVGName,
		"IconSVG":    body.IconSVG,
		"Hyperlink":  body.Hyperlink,
		"LinkName":   body.LinkName,
		"Category":   body.Category,
		"Enabled":    body.Enabled,
		"ModifiedBy": username,
	}

	_, err2 := db.ExecuteStoredProcedure("dbo.PR_ExternalLinks_Update", params)
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

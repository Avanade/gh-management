package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	session "main/pkg/session"
	// "main/pkg/sql"

	"net/http"
	// "os"
	ghmgmt "main/pkg/ghmgmtdb"
	"github.com/gorilla/mux"
)



func GetExternalLinks(w http.ResponseWriter, r *http.Request) {

	ExternalLinks, err := ghmgmt.ExternalLinksExecuteSelect()
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

func GetExternalLinksAllEnabled(w http.ResponseWriter, r *http.Request) {
	param := map[string]interface{}{
		"Enabled": true,
	}

	ExternalLinks, err := ghmgmt.ExternalLinksExecuteAllEnabled(param)
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

	param := map[string]interface{}{
		"Id": id,
	}

	ExternalLinks, err := ghmgmt.ExternalLinksExecuteById(param)
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

	params := map[string]interface{}{
		"IconSVG":   data.IconSVG,
		"Hyperlink": data.Hyperlink,
		"LinkName":  data.LinkName,
		"Enabled":   data.Enabled,
		"CreatedBy": username,
	}

	__, err := ghmgmt.ExternalLinksExecuteCreate(params)
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

	params := map[string]interface{}{
		"Id":         body.Id,
		"IconSVG":    body.IconSVG,
		"Hyperlink":  body.Hyperlink,
		"LinkName":   body.LinkName,
		"Enabled":    body.Enabled,
		"ModifiedBy": username,
	}

	_, err2 := ghmgmt.ExternalLinksExecuteUpdate(params)
	if err != nil {
		fmt.Println(err2)
		return
	}
}

func ExternalLinksDelete(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	// cp := sql.ConnectionParam{
	// 	ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	// }
	// db, err := sql.Init(cp)
	param := map[string]interface{}{
		"Id": id,
	}
	_, error := ghmgmt.ExternalLinksExecuteDelete(param)
	if error != nil {
		fmt.Println(error)
	}
}
package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"main/models"
	db "main/pkg/ghmgmtdb"
	"main/pkg/session"

	"github.com/gorilla/mux"
)

func GetExternalLinks(w http.ResponseWriter, r *http.Request) {

	ExternalLinks, err := db.ExternalLinksExecuteSelect()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(ExternalLinks)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetExternalLinksAllEnabled(w http.ResponseWriter, r *http.Request) {
	param := map[string]interface{}{
		"Enabled": true,
	}

	ExternalLinks, err := db.ExternalLinksExecuteAllEnabled(param)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(ExternalLinks)
	if err != nil {
		log.Println(err.Error())
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

	ExternalLinks, err := db.ExternalLinksExecuteById(param)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(ExternalLinks)
	if err != nil {
		log.Println(err.Error())
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
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := map[string]interface{}{
		"IconSVG":   data.IconSVG,
		"Hyperlink": data.Hyperlink,
		"LinkName":  data.LinkName,
		"Enabled":   data.Enabled,
		"CreatedBy": username,
	}

	_, err = db.ExternalLinksExecuteCreate(params)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	_, err = db.ExternalLinksExecuteUpdate(params)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ExternalLinksDelete(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	param := map[string]interface{}{
		"Id": id,
	}
	_, err := db.ExternalLinksExecuteDelete(param)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

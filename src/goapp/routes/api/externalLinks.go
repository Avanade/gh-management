package routes

import (
	"encoding/json"
	"net/http"

	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"
	"main/pkg/session"

	"github.com/gorilla/mux"
)

type ExternalLinksDto struct {
	Id int `json:"id"`

	IconSVG   string `json:"iconsvg"`
	Hyperlink string `json:"hyperlink"`
	LinkName  string `json:"linkname"`

	Enabled    string `json:"enabled"`
	Created    string `json:"created"`
	CreatedBy  string `json:"createdBy"`
	Modified   string `json:"modified"`
	ModifiedBy string `json:"modifiedBy"`
}

func GetExternalLinks(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	externalLinks, err := db.ExternalLinksExecuteSelect()
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(externalLinks)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetExternalLinksEnabled(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	param := map[string]interface{}{
		"IsEnabled": true,
	}

	externalLinks, err := db.ExternalLinksExecuteAllEnabled(param)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(externalLinks)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetExternalLinkById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id := req["id"]

	param := map[string]interface{}{
		"Id": id,
	}

	externalLinks, err := db.ExternalLinksExecuteById(param)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(externalLinks)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func CreateExternalLinks(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	var data ExternalLinksDto

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := map[string]interface{}{
		"IconSVG":   data.IconSVG,
		"Hyperlink": data.Hyperlink,
		"LinkName":  data.LinkName,
		"IsEnabled": data.Enabled,
		"CreatedBy": username,
	}

	_, err = db.ExternalLinksExecuteCreate(params)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateExternalLinksById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	var body ExternalLinksDto

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := map[string]interface{}{
		"Id":         body.Id,
		"IconSVG":    body.IconSVG,
		"Hyperlink":  body.Hyperlink,
		"LinkName":   body.LinkName,
		"IsEnabled":  body.Enabled,
		"ModifiedBy": username,
	}

	_, err = db.ExternalLinksExecuteUpdate(params)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteExternalLinkById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id := req["id"]

	param := map[string]interface{}{
		"Id": id,
	}
	_, err := db.ExternalLinksExecuteDelete(param)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

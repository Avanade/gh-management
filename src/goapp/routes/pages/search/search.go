package routes

import (
	"encoding/json"
	"log"
	"net/http"

	db "main/pkg/ghmgmtdb"
	"main/pkg/session"
	"main/pkg/template"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	template.UseTemplate(&w, r, "search/search", nil)
}

func GetSearchResults(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	search := params.Get("search")
	offset := params.Get("offset")
	filter := params.Get("filter")
	selectedSourceType := params.Get("selectedSourceType")

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"].(string)

	data, total, err := db.SearchCommunitiesProjectsUsers(search, offset, filter, selectedSourceType, username)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(struct {
		Data  []map[string]interface{} `json:"data"`
		Total int                      `json:"total"`
	}{
		Data:  data,
		Total: total,
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

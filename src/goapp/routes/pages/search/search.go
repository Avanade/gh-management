package routes

import (
	"encoding/json"
	"log"
	"net/http"

	db "main/pkg/ghmgmtdb"
	"main/pkg/session"
	"main/pkg/template"

	"github.com/gorilla/mux"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	template.UseTemplate(&w, r, "search/search", nil)

}

func GetSearchResults(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)

	searchText := r.URL.Query().Get("search")
	offSet := req["offSet"]
	rowCount := req["rowCount"]
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"].(string)

	searchResults, err := db.SearchCommunitiesProjectsUsers(searchText, offSet, rowCount, username)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(searchResults)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

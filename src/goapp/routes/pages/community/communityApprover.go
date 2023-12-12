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

type CommunityApprovers struct {
	Id                        int    `json:"id"`
	ApproverUserPrincipalName string `json:"name"`
	Disabled                  bool   `json:"disabled"`
	Created                   string `json:"created"`
	CreatedBy                 string `json:"createdBy"`
	Modified                  string `json:"modified"`
	ModifiedBy                string `json:"modifiedBy"`
}

func CommunityApproversHandler(w http.ResponseWriter, r *http.Request) {

	template.UseTemplate(&w, r, "/community/communityapprovers", nil)
}

func GetCommunityApproversById(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	approvers, err := db.GetCommunityApproversById(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvers)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetCommunityApproversList(w http.ResponseWriter, r *http.Request) {
	approvers, err := db.GetCommunityApprovers()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvers)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetAllActiveCommunityApprovers(w http.ResponseWriter, r *http.Request) {
	approvers, err := db.GetActiveCommunityApprovers()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvers)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func CommunityApproversListUpdate(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"].(string)
	var body CommunityApprovers

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.UpdateCommunityApproversById(body.Id, body.Disabled, body.ApproverUserPrincipalName, username)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

package routes

import (
	"encoding/json"

	session "main/pkg/session"
	template "main/pkg/template"
	"net/http"

	"fmt"
	models "main/models"

	db "main/pkg/ghmgmtdb"

	"github.com/gorilla/mux"
)

func CommunityApproverHandler(w http.ResponseWriter, r *http.Request) {

	template.UseTemplate(&w, r, "/community/communityapprovers", nil)
}

func GetCommunityApproversById(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	approvers, err := db.GetCommunityApproversById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetCommunityApproversList(w http.ResponseWriter, r *http.Request) {
	approvers, err := db.GetCommunityApprovers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetAllActiveCommunityApprovers(w http.ResponseWriter, r *http.Request) {
	approvers, err := db.GetActiveCommunityApprovers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(approvers)
	if err != nil {
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
	var body models.TypCommunityApprovers

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	_, errDb := db.UpdateCommunityApproversById(body.Id, body.Disabled, body.ApproverUserPrincipalName, username)
	if errDb != nil {
		if errDb != nil {
			http.Error(w, errDb.Error(), http.StatusBadRequest)
			fmt.Println(errDb)
			return
		}
	}
}

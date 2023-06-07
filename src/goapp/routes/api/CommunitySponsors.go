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

func CommunitySponsorsAPIHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	var body models.TypCommunitySponsors
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
		param := map[string]interface{}{

			"CommunityId":       body.CommunityId,
			"UserPrincipalName": body.UserPrincipalName,
			"CreatedBy":         username,
		}

		_, err := db.CommunitySponsorsInsert(param)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "PUT":
		param := map[string]interface{}{
			"CommunityId":       body.CommunityId,
			"UserPrincipalName": body.UserPrincipalName,
			"CreatedBy":         username,
		}
		_, err := db.CommunitySponsorsUpdate(param)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CommunitySponsorsPerCommunityId(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	// Get project list
	param := map[string]interface{}{

		"CommunityId": id,
	}

	CommunitySponsors, err := db.CommunitySponsorsSelectByCommunityId(param)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(CommunitySponsors)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

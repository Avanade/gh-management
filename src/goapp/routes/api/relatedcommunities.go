package routes

import (
	"encoding/json"
	"log"
	models "main/models"
	ghmgmt "main/pkg/ghmgmtdb"
	"net/http"

	"github.com/gorilla/mux"
)

func RelatedCommunitiesInsert(w http.ResponseWriter, r *http.Request) {

	var body models.TypRelatedCommunity

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	param := map[string]interface{}{

		"ParentCommunityId":  body.ParentCommunityId,
		"RelatedCommunityId": body.RelatedCommunityId,
	}

	approvers, err := ghmgmt.RelatedCommunitiesInsert(param)
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

func RelatedCommunitiesDelete(w http.ResponseWriter, r *http.Request) {

	var body models.TypRelatedCommunity

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	param := map[string]interface{}{

		"ParentCommunityId":  body.ParentCommunityId,
		"RelatedCommunityId": body.RelatedCommunityId,
	}

	approvers, err := ghmgmt.RelatedCommunitiesDelete(param)
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

func RelatedCommunitiesSelect(w http.ResponseWriter, r *http.Request) {

	req := mux.Vars(r)
	id := req["id"]

	param := map[string]interface{}{

		"ParentCommunityId": id,
	}

	approvers, err := ghmgmt.RelatedCommunitiesSelect(param)
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

package routes

import (
	"encoding/json"
	"log"
	ghmgmt "main/pkg/ghmgmtdb"
	"net/http"

	"github.com/gorilla/mux"
)

func CommunityTagPerCommunityId(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	// Get project list
	param := map[string]interface{}{

		"CommunityId": id,
	}

	CommunityTags, err := ghmgmt.CommunityTagsSelectByCommunityId(param)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(CommunityTags)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

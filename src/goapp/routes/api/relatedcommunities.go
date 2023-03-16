package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	"main/pkg/sql"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func RelatedCommunitiesInsert(w http.ResponseWriter, r *http.Request) {

	var body models.TypRelatedCommunity

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)

		return
	}

	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)
	defer db.Close()

	param := map[string]interface{}{

		"ParentCommunityId":  body.ParentCommunityId,
		"RelatedCommunityId": body.RelatedCommunityId,
	}

	approvers, err := db.ExecuteStoredProcedure("PR_RelatedCommunities_Insert", param)
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

func RelatedCommunitiesDelete(w http.ResponseWriter, r *http.Request) {

	var body models.TypRelatedCommunity

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)

		return
	}

	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)
	defer db.Close()

	param := map[string]interface{}{

		"ParentCommunityId":  body.ParentCommunityId,
		"RelatedCommunityId": body.RelatedCommunityId,
	}

	approvers, err := db.ExecuteStoredProcedure("PR_RelatedCommunities_Delete", param)
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

func RelatedCommunitiesSelect(w http.ResponseWriter, r *http.Request) {

	req := mux.Vars(r)
	id := req["id"]

	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)
	defer db.Close()

	param := map[string]interface{}{

		"ParentCommunityId": id,
	}

	approvers, err := db.ExecuteStoredProcedureWithResult("PR_RelatedCommunities_Select", param)
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

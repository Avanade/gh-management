package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	"main/pkg/sql"
	"net/http"
	"os"
)

func RelatedCommunitiesInsert(w http.ResponseWriter, r *http.Request) {

	var body models.TypRelatedCommunities

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

	approvers, err := db.ExecuteStoredProcedureWithResult("PR_RelatedCommunities_Insert", param)
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

	var body models.TypRelatedCommunities

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

	approvers, err := db.ExecuteStoredProcedureWithResult("PR_RelatedCommunities_Insert", param)
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

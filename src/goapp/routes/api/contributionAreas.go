package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"

	"github.com/gorilla/mux"
)

type ContributionAreaDto struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Created    string `json:"created"`
	CreatedBy  string `json:"createdBy"`
	Modified   string `json:"modified"`
	ModifiedBy string `json:"modifiedBy"`
}

func GetContributionAreasByActivityId(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	vars := mux.Vars(r)
	activityId, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := db.AdditionalContributionAreas_Select(activityId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

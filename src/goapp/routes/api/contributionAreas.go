package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"
	"main/pkg/session"

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

func GetContributionAreas(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	var data interface{}
	var total int

	params := r.URL.Query()

	if params.Has("offset") && params.Has("filter") {
		filter, _ := strconv.Atoi(params["filter"][0])
		offset, _ := strconv.Atoi(params["offset"][0])
		search := params["search"][0]
		orderby := params["orderby"][0]
		ordertype := params["ordertype"][0]
		data, _ = db.ContributionAreas_SelectByFilter(offset, filter, orderby, ordertype, search)
	} else {
		result, err := db.ContributionAreas_Select()
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = result
	}

	total = db.SelectTotalContributionAreas()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Data  interface{} `json:"data"`
		Total int         `json:"total"`
	}{
		Data:  data,
		Total: total,
	})
}

func GetContributionAreaById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := db.GetContributionAreaById(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
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

func CreateContributionAreas(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var contributionArea ContributionAreaDto
	json.NewDecoder(r.Body).Decode(&contributionArea)

	id, err := db.ContributionAreas_Insert(contributionArea.Name, username)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	contributionArea.Id = id
	json.NewEncoder(w).Encode(contributionArea)
}

func UpdateContributionArea(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var contributionArea ContributionAreaDto
	json.NewDecoder(r.Body).Decode(&contributionArea)

	db.UpdateContributionAreaById(contributionArea.Id, contributionArea.Name, username)
	json.NewEncoder(w).Encode(contributionArea)
}

package routes

import (
	"encoding/json"
	"net/http"

	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"
)

type OssContributionSponsor struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	IsArchived bool   `json:"isArchived"`
}

func GetAllOssContributionSponsors(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sponsors, err := db.SelectAllSponsors()
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(sponsors)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func GetAllEnabledOssContributionSponsors(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	param := map[string]interface{}{
		"IsArchived": false,
	}

	sponsors, err := db.SelectSponsorsByIsArchived(param)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(sponsors)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func AddSponsor(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	var data OssContributionSponsor

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := map[string]interface{}{
		"Name":       data.Name,
		"IsArchived": data.IsArchived,
	}

	_, err = db.InsertSponsor(params)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateSponsor(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	var data OssContributionSponsor

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := map[string]interface{}{
		"Id":         data.Id,
		"Name":       data.Name,
		"IsArchived": data.IsArchived,
	}

	_, err = db.UpdateSponsor(params)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

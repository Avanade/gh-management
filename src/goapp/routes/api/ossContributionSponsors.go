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

func MigrateToOssSponsorsTable(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sponsorsInitList := []string{
		"Solution Area - Applications and Infrastructure",
		"Solution Area - Business Applications",
		"Solution Area - Data & AI",
		"Solution Area - Modern Workplace",
		"Business Group - Advisory",
		"ITS",
		"Office of the CTO",
	}

	for _, sponsor := range sponsorsInitList {
		param := map[string]interface{}{
			"Name": sponsor,
		}

		result, err := db.SelectSponsorByName(param)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(result) == 0 {
			params := map[string]interface{}{
				"Name":       sponsor,
				"IsArchived": false,
			}

			_, err = db.InsertSponsor(params)
			if err != nil {
				logger.LogException(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	repos, err := db.SelectReposWithMakePublicRequest()
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, repo := range repos {

		params := map[string]interface{}{
			"Name": repo["OSSsponsor"].(string),
		}

		sponsors, err := db.SelectSponsorByName(params)
		if err != nil {
			continue
		}

		for _, sponsor := range sponsors {
			params2 := map[string]interface{}{
				"Id":                       repo["Id"],
				"OSSContributionSponsorId": sponsor["Id"],
			}

			_, err := db.UpdateOssContributionSponsorId(params2)
			if err != nil {
				continue
			}
		}

	}
}

package routes

import (
	"encoding/json"
	"net/http"

	db "main/pkg/ghmgmtdb"
)

type ContributionAreaDto struct {
	Id         int    `json:id`
	Name       string `json:name`
	Created    string `json:created`
	CreatedBy  string `json:createdBy`
	Modified   string `json:modified`
	ModifiedBy string `json:modifiedBy`
}

func GetContributionAreas(w http.ResponseWriter, r *http.Request) {
	result := db.ContributionAreas_Select()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func CreateContributionAreas(w http.ResponseWriter, r *http.Request) {
	var contributionArea ContributionAreaDto
	json.NewDecoder(r.Body).Decode(&contributionArea)

	id, err := db.ContributionAreas_Insert(contributionArea.Name, contributionArea.CreatedBy)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	contributionArea.Id = id
	json.NewEncoder(w).Encode(contributionArea)
}

package routes

import (
	"encoding/json"
	"net/http"

	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"
)

type ActivityTypeDto struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetActivityTypes(w http.ResponseWriter, r *http.Request) {
	result := db.ActivityTypes_Select()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func CreateActivityType(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	var activityType ActivityTypeDto
	json.NewDecoder(r.Body).Decode(&activityType)
	id, err := db.ActivityTypes_Insert(activityType.Name)
	if err != nil {
		logger.LogException(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	activityType.Id = id
	json.NewEncoder(w).Encode(activityType)
}

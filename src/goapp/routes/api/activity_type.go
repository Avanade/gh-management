package routes

import (
	"encoding/json"
	"net/http"

	db "main/pkg/ghmgmtdb"
)

type ActivityTypeDto struct {
	Id   int    `json:id`
	Name string `json:name`
}

func GetActivityTypes(w http.ResponseWriter, r *http.Request) {
	result := db.PRActivityTypes_Select()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func CreateActivityType(w http.ResponseWriter, r *http.Request) {
	var activityType ActivityTypeDto
	json.NewDecoder(r.Body).Decode(&activityType)
	id, err := db.PRActivityTypes_Insert(activityType.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	activityType.Id = id
	json.NewEncoder(w).Encode(activityType)
}

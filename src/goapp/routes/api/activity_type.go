package routes

import (
	"encoding/json"
	"log"
	"net/http"

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
	var activityType ActivityTypeDto
	json.NewDecoder(r.Body).Decode(&activityType)
	id, err := db.ActivityTypes_Insert(activityType.Name)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	activityType.Id = id
	json.NewEncoder(w).Encode(activityType)
}

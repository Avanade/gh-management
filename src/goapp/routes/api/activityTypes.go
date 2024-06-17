package routes

import (
	"encoding/json"
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

package routes

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ApprovalTypeDto struct {
	Id                        int    `json:id`
	Name                      string `json:name`
	ApproverUserPrincipalName string `json:approver_user_principal_name`
	IsActive                  bool   `json:is_active`
}

func GetApprovalTypes(w http.ResponseWriter, r *http.Request) {
	// result := db.PRActivityTypes_Select()
	result := ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func GetApprovalTypeById(w http.ResponseWriter, r *http.Request) {
	// result := db.PRActivityTypes_Select()
	result := ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func CreateApprovalType(w http.ResponseWriter, r *http.Request) {
	var approvalTypeDto ApprovalTypeDto
	json.NewDecoder(r.Body).Decode(&approvalTypeDto)
	// id, err := db.PRActivityTypes_Insert(activityType.Name)
	id, err := 0, errors.New("THIS IS ERROR")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	approvalTypeDto.Id = id
	json.NewEncoder(w).Encode(approvalTypeDto)
}

func EditApprovalTypeById(w http.ResponseWriter, r *http.Request) {
	var approvalTypeDto ApprovalTypeDto
	json.NewDecoder(r.Body).Decode(&approvalTypeDto)
	// id, err := db.PRActivityTypes_Insert(activityType.Name)
	id, err := 0, errors.New("THIS IS ERROR")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	approvalTypeDto.Id = id
	json.NewEncoder(w).Encode(approvalTypeDto)
}

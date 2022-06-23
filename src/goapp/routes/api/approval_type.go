package routes

import (
	"encoding/json"
	"fmt"
	db "main/pkg/ghmgmtdb"
	session "main/pkg/session"
	"net/http"
	"strconv"

	"main/models"

	"github.com/gorilla/mux"
)

type ApprovalTypeDto struct {
	Id                        int    `json:id`
	Name                      string `json:name`
	ApproverUserPrincipalName string `json:approver_user_principal_name`
	IsActive                  bool   `json:is_active`
}

func GetApprovalTypes(w http.ResponseWriter, r *http.Request) {
	result, err := db.SelectApprovalTypes()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func GetApprovalTypeById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	result, err := db.SelectApprovalTypeById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func CreateApprovalType(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var approvalTypeDto ApprovalTypeDto
	json.NewDecoder(r.Body).Decode(&approvalTypeDto)
	id, err := db.InsertApprovalType(models.ApprovalType{
		Name:                      approvalTypeDto.Name,
		ApproverUserPrincipalName: approvalTypeDto.ApproverUserPrincipalName,
		IsActive:                  approvalTypeDto.IsActive,
		CreatedBy:                 username,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	approvalTypeDto.Id = id
	json.NewEncoder(w).Encode(approvalTypeDto)
}

func EditApprovalTypeById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var approvalTypeDto ApprovalTypeDto
	json.NewDecoder(r.Body).Decode(&approvalTypeDto)

	id, _ := strconv.Atoi(vars["id"])
	approvalTypeId, err := db.UpdateApprovalType(models.ApprovalType{
		Id:                        id,
		Name:                      approvalTypeDto.Name,
		ApproverUserPrincipalName: approvalTypeDto.ApproverUserPrincipalName,
		IsActive:                  approvalTypeDto.IsActive,
		CreatedBy:                 username,
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	approvalTypeDto.Id = approvalTypeId
	json.NewEncoder(w).Encode(approvalTypeDto)
}

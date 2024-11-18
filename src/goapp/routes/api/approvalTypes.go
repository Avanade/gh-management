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

type ApprovalTypeDto struct {
	Id         int           `json:"id"`
	Name       string        `json:"name"`
	Approvers  []ApproverDto `json:"approvers"`
	IsActive   bool          `json:"isActive"`
	IsArchived bool          `json:"isArchived"`
}

type ApproverDto struct {
	ApprovalTypeId int    `json:"approvalTypeId"`
	ApproverEmail  string `json:"approverEmail"`
	ApproverName   string `json:"approverName"`
}

func EditApprovalTypeById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	vars := mux.Vars(r)

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var approvalTypeDto ApprovalTypeDto
	json.NewDecoder(r.Body).Decode(&approvalTypeDto)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	approvalTypeId, err := db.UpdateApprovalType(db.ApprovalType{
		Id:        id,
		Name:      approvalTypeDto.Name,
		IsActive:  approvalTypeDto.IsActive,
		CreatedBy: username,
	})
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.DeleteApproverByApprovalTypeId(id)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range approvalTypeDto.Approvers {
		err = db.InsertUser(v.ApproverEmail, v.ApproverName, "", "", "")
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := db.InsertApprover(db.Approver{
			ApprovalTypeId: id,
			ApproverEmail:  v.ApproverEmail,
		})
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	approvalTypeDto.Id = approvalTypeId
	json.NewEncoder(w).Encode(approvalTypeDto)
}

func SetIsArchivedApprovalTypeById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	vars := mux.Vars(r)

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var approvalTypeDto ApprovalTypeDto
	json.NewDecoder(r.Body).Decode(&approvalTypeDto)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	approvalTypeId, _, err := db.SetIsArchiveApprovalTypeById(db.ApprovalType{
		Id:         id,
		Name:       approvalTypeDto.Name,
		IsArchived: approvalTypeDto.IsArchived,
		ModifiedBy: username,
	})

	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	approvalTypeDto.Id = approvalTypeId
	json.NewEncoder(w).Encode(approvalTypeDto)
}

func GetActiveApprovalTypes(w http.ResponseWriter, r *http.Request) {

	data := db.GetAllActiveApprovers()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func getApproversByApprovalTypeId(approvalTypeId int) (*[]ApproverDto, error) {
	resultApprovers, err := db.GetApproversByApprovalTypeId(approvalTypeId)
	if err != nil {
		return nil, err
	}

	var approversDto []ApproverDto

	for _, v := range resultApprovers {
		approverDto := ApproverDto{
			ApprovalTypeId: v.ApprovalTypeId,
			ApproverEmail:  v.ApproverEmail,
			ApproverName:   v.ApproverName,
		}

		approversDto = append(approversDto, approverDto)
	}

	return &approversDto, nil
}

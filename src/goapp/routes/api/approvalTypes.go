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

func GetApprovalTypes(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	var data []map[string]interface{}
	var total int

	params := r.URL.Query()

	if params.Has("offset") && params.Has("filter") {
		filter, _ := strconv.Atoi(params["filter"][0])
		offset, _ := strconv.Atoi(params["offset"][0])
		search := params["search"][0]
		orderby := params["orderby"][0]
		ordertype := params["ordertype"][0]
		result, err := db.SelectApprovalTypesByFilter(offset, filter, orderby, ordertype, search)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data = result
	} else {
		result, err := db.SelectApprovalTypes()
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data = result
	}

	//MOVE APPROVAL TYPES FROM DATABASE RESULT TO API DTO
	var approvalTypesDto []ApprovalTypeDto
	for _, v := range data {
		approvalTypeDto := ApprovalTypeDto{
			Id:         int(v["Id"].(int64)),
			Name:       v["Name"].(string),
			IsActive:   v["IsActive"].(bool),
			IsArchived: v["IsArchived"].(bool),
		}

		approversResult, err := getApproversByApprovalTypeId(approvalTypeDto.Id)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		approvalTypeDto.Approvers = *approversResult

		approvalTypesDto = append(approvalTypesDto, approvalTypeDto)
	}

	total = db.SelectTotalApprovalTypes()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Data  []ApprovalTypeDto `json:"data"`
		Total int               `json:"total"`
	}{
		Data:  approvalTypesDto,
		Total: total,
	})
}

func GetApprovalTypeById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	result, err := db.SelectApprovalTypeById(id)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	approversDto, err := getApproversByApprovalTypeId(result.Id)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	approvalTypeDto := ApprovalTypeDto{
		Id:         result.Id,
		Name:       result.Name,
		Approvers:  *approversDto,
		IsActive:   result.IsActive,
		IsArchived: result.IsArchived,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(approvalTypeDto)
}

func CreateApprovalType(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var approvalTypeDto ApprovalTypeDto
	json.NewDecoder(r.Body).Decode(&approvalTypeDto)
	id, err := db.InsertApprovalType(db.ApprovalType{
		Name:      approvalTypeDto.Name,
		IsActive:  approvalTypeDto.IsActive,
		CreatedBy: username,
	})
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range approvalTypeDto.Approvers {
		err = db.InsertUser(v.ApproverEmail, v.ApproverName, "", "", "")
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

	approvalTypeDto.Id = id
	json.NewEncoder(w).Encode(approvalTypeDto)
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

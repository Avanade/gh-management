package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	db "main/pkg/ghmgmtdb"
	"main/pkg/session"

	"github.com/gorilla/mux"
)

type ApprovalTypeDto struct {
	Id                        int           `json:"id"`
	Name                      string        `json:"name"`
	ApproverUserPrincipalName string        `json:"approverUserPrincipalName"`
	Approvers                 []ApproverDto `json:"approvers"`
	IsActive                  bool          `json:"isActive"`
	IsArchived                bool          `json:"isArchived"`
}

type ApproverDto struct {
	ApprovalTypeId int    `json:"approvalTypeId"`
	ApproverEmail  string `json:"approverEmail"`
}

func GetApprovalTypes(w http.ResponseWriter, r *http.Request) {
	var data interface{}
	var total int

	params := r.URL.Query()

	if params.Has("offset") && params.Has("filter") {
		filter, _ := strconv.Atoi(params["filter"][0])
		offset, _ := strconv.Atoi(params["offset"][0])
		search := params["search"][0]
		orderby := params["orderby"][0]
		ordertype := params["ordertype"][0]
		data, _ = db.SelectApprovalTypesByFilter(offset, filter, orderby, ordertype, search)
	} else {
		result, err := db.SelectApprovalTypes()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data = result
	}

	total = db.SelectTotalApprovalTypes()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Data  interface{} `json:"data"`
		Total int         `json:"total"`
	}{
		Data:  data,
		Total: total,
	})
}

func GetApprovalTypeById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	result, err := db.SelectApprovalTypeById(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultApprovers, err := db.GetApproversByApprovalTypeId(result.Id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var approversDto []ApproverDto

	for _, v := range resultApprovers {
		approverDto := ApproverDto{
			ApprovalTypeId: v.ApprovalTypeId,
			ApproverEmail:  v.ApproverEmail,
		}

		approversDto = append(approversDto, approverDto)
	}

	approvalTypeDto := ApprovalTypeDto{
		Id:                        result.Id,
		Name:                      result.Name,
		ApproverUserPrincipalName: result.ApproverUserPrincipalName,
		Approvers:                 approversDto,
		IsActive:                  result.IsActive,
		IsArchived:                result.IsArchived,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(approvalTypeDto)
}

func CreateApprovalType(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var approvalTypeDto ApprovalTypeDto
	json.NewDecoder(r.Body).Decode(&approvalTypeDto)
	id, err := db.InsertApprovalType(db.ApprovalType{
		Name:                      approvalTypeDto.Name,
		ApproverUserPrincipalName: approvalTypeDto.ApproverUserPrincipalName,
		IsActive:                  approvalTypeDto.IsActive,
		CreatedBy:                 username,
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range approvalTypeDto.Approvers {
		err := db.InsertApprover(db.Approver{
			ApprovalTypeId: id,
			ApproverEmail:  v.ApproverEmail,
		})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	approvalTypeId, err := db.UpdateApprovalType(db.ApprovalType{
		Id:                        id,
		Name:                      approvalTypeDto.Name,
		ApproverUserPrincipalName: approvalTypeDto.ApproverUserPrincipalName,
		IsActive:                  approvalTypeDto.IsActive,
		CreatedBy:                 username,
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.DeleteApproverByApprovalTypeId(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range approvalTypeDto.Approvers {
		err := db.InsertApprover(db.Approver{
			ApprovalTypeId: id,
			ApproverEmail:  v.ApproverEmail,
		})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	approvalTypeDto.Id = approvalTypeId
	json.NewEncoder(w).Encode(approvalTypeDto)
}

func SetIsArchivedApprovalTypeById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var approvalTypeDto ApprovalTypeDto
	json.NewDecoder(r.Body).Decode(&approvalTypeDto)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	approvalTypeId, _, err := db.SetIsArchiveApprovalTypeById(db.ApprovalType{
		Id:                        id,
		Name:                      approvalTypeDto.Name,
		ApproverUserPrincipalName: approvalTypeDto.ApproverUserPrincipalName,
		IsArchived:                approvalTypeDto.IsArchived,
		ModifiedBy:                username,
	})

	if err != nil {
		log.Println(err.Error())
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

package approvaltype

import (
	"encoding/json"
	"fmt"
	"main/model"
	"main/pkg/appinsights_wrapper"
	"main/pkg/session"
	"main/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type approvalTypeController struct {
	*service.Service
}

func NewApprovalTypeController(service *service.Service) ApprovalTypeController {
	return &approvalTypeController{service}
}

func (c *approvalTypeController) CreateApprovalType(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var approvalType model.ApprovalType
	json.NewDecoder(r.Body).Decode(&approvalType)
	approvalType.CreatedBy = username
	id, err := c.ApprovalType.Insert(&approvalType)
	if err != nil {
		logger.TrackException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range approvalType.Approvers {
		err = c.User.Create(&model.User{
			UserPrincipalName: v.ApproverEmail,
			Name:              v.ApproverName,
			GivenName:         "",
			Surname:           "",
			JobTitle:          "",
		})
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v.ApprovalTypeId = id
		err := c.RepositoryApprover.Create(&v)
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	json.NewEncoder(w).Encode(id)
}

func (c *approvalTypeController) GetApprovalTypeById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	approvalType, err := c.ApprovalType.GetById(id)
	if err != nil {
		logger.TrackException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(approvalType)
}

func (c *approvalTypeController) GetApprovalTypes(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	var approvalTypes []model.ApprovalType
	var total int64
	var err error

	params := r.URL.Query()
	if params.Has("offset") && params.Has("filter") {
		filter, _ := strconv.Atoi(params["filter"][0])
		offset, _ := strconv.Atoi(params["offset"][0])
		opt := model.FilterOptions{
			Filter:    filter,
			Offset:    offset,
			Search:    params.Get("search"),
			Orderby:   params.Get("orderby"),
			Ordertype: params.Get("ordertype"),
		}
		approvalTypes, total, err = c.ApprovalType.Get(&opt)
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		approvalTypes, total, err = c.ApprovalType.Get(nil)
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	for i, v := range approvalTypes {
		approvers, err := c.RepositoryApprover.Get(v.Id)
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		approvalTypes[i].Approvers = approvers
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Data  []model.ApprovalType `json:"data"`
		Total int64                `json:"total"`
	}{
		Data:  approvalTypes,
		Total: total,
	})
}

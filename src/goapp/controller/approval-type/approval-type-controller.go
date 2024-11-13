package approvaltype

import (
	"encoding/json"
	"main/model"
	"main/pkg/appinsights_wrapper"
	"main/service"
	"net/http"
	"strconv"
)

type approvalTypeController struct {
	*service.Service
}

func NewApprovalTypeController(service *service.Service) ApprovalTypeController {
	return &approvalTypeController{service}
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
		approvers, err := c.Approver.Get(v.Id)
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

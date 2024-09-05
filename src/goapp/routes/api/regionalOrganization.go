package routes

import (
	"encoding/json"
	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type RegionalOrganizationDto struct {
	Id                      int64     `json:"id"`
	Name                    string    `json:"name"`
	IsCleanUpMembersEnabled bool      `json:"isCleanUpMembersEnabled"`
	IsIndexRepoEnabled      bool      `json:"isIndexRepoEnabled"`
	IsCopilotRequestEnabled bool      `json:"isCopilotRequestEnabled"`
	IsAccessRequestEnabled  bool      `json:"isAccessRequestEnabled"`
	IsEnabled               bool      `json:"isEnabled"`
	Created                 time.Time `json:"created"`
	CreatedBy               string    `json:"createdBy"`
	ModifiedBy              string    `json:"modifiedBy"`
	Modified                time.Time `json:"modified"`
}

func GetRegionalOrganizations(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	regionalOrganizations, err := db.SelectRegionalOrganization()
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(regionalOrganizations)
}

func GetRegionalOrganizationByOption(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	query := r.URL.Query()
	var offset int64 = 0
	var filter int64 = 10
	if query.Get("offset") != "" {
		offset, _ = strconv.ParseInt(query.Get("offset"), 10, 64)
	}
	if query.Get("filter") != "" {
		filter, _ = strconv.ParseInt(query.Get("filter"), 10, 64)
	}
	search := query.Get("search")
	orderBy := query.Get("orderBy")
	orderType := query.Get("orderType")

	regionalOrganizations, total, err := db.SelectRegionalOrganizationByOption(offset, filter, search, orderBy, orderType)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var regionalOrganizationsDto []RegionalOrganizationDto
	for _, regionalOrganization := range regionalOrganizations {
		regionalOrganizationDto := RegionalOrganizationDto{
			Id:                      regionalOrganization.Id,
			Name:                    regionalOrganization.Name,
			IsCleanUpMembersEnabled: regionalOrganization.IsCleanUpMembersEnabled,
			IsIndexRepoEnabled:      regionalOrganization.IsIndexRepoEnabled,
			IsCopilotRequestEnabled: regionalOrganization.IsCopilotRequestEnabled,
			IsAccessRequestEnabled:  regionalOrganization.IsAccessRequestEnabled,
			IsEnabled:               regionalOrganization.IsEnabled,
		}
		regionalOrganizationsDto = append(regionalOrganizationsDto, regionalOrganizationDto)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", strconv.FormatInt(total, 10))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Data  []RegionalOrganizationDto `json:"data"`
		Total int64                     `json:"total"`
	}{
		Data:  regionalOrganizationsDto,
		Total: total,
	})
}

func GetRegionalOrganizationById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id, err := strconv.ParseInt(req["id"], 10, 64)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	regionalOrganization, err := db.SelectRegionalOrganizationById(id)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(regionalOrganization)
}

func InsertRegionalOrganization(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	var regionalOrganizationDto RegionalOrganizationDto
	err := json.NewDecoder(r.Body).Decode(&regionalOrganizationDto)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	regionalOrganization := db.RegionalOrganization{
		Id:                      regionalOrganizationDto.Id,
		Name:                    regionalOrganizationDto.Name,
		IsCleanUpMembersEnabled: regionalOrganizationDto.IsCleanUpMembersEnabled,
		IsIndexRepoEnabled:      regionalOrganizationDto.IsIndexRepoEnabled,
		IsCopilotRequestEnabled: regionalOrganizationDto.IsCopilotRequestEnabled,
		IsAccessRequestEnabled:  regionalOrganizationDto.IsAccessRequestEnabled,
		IsEnabled:               regionalOrganizationDto.IsEnabled,
	}

	_, err = db.InsertRegionalOrganization(regionalOrganization)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(regionalOrganizationDto)
}

func UpdateRegionalOrganization(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id, err := strconv.ParseInt(req["id"], 10, 64)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var regionalOrganizationDto RegionalOrganizationDto
	err = json.NewDecoder(r.Body).Decode(&regionalOrganizationDto)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	regionalOrganization := db.RegionalOrganization{
		Id:                      id,
		Name:                    regionalOrganizationDto.Name,
		IsCleanUpMembersEnabled: regionalOrganizationDto.IsCleanUpMembersEnabled,
		IsIndexRepoEnabled:      regionalOrganizationDto.IsIndexRepoEnabled,
		IsCopilotRequestEnabled: regionalOrganizationDto.IsCopilotRequestEnabled,
		IsAccessRequestEnabled:  regionalOrganizationDto.IsAccessRequestEnabled,
		IsEnabled:               regionalOrganizationDto.IsEnabled,
	}

	err = db.UpdateRegionalOrganization(regionalOrganization)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	regionalOrganizationDto.Id = id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(regionalOrganizationDto)
}

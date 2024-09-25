package routes

import (
	"encoding/json"
	"fmt"
	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"
	ghAPI "main/pkg/github"
	"main/pkg/session"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type RegionalOrganizationDto struct {
	Id                      int64     `json:"id"`
	Name                    string    `json:"name"`
	IsRegionalOrganization  bool      `json:"isRegionalOrganization"`
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

type EnterpriseOrganization struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func GetEnterpriseOrganizations(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	token := os.Getenv("GH_ENTERPRISE_TOKEN")
	enterprise := os.Getenv("GH_ENTERPRISE_NAME")
	enterpriseOrgs, err := ghAPI.GetOrganizationsWithinEnterprise(enterprise, token)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isEnabled := db.NullBool{Value: true}
	regionalOrganizations, err := db.SelectRegionalOrganization(&isEnabled)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filteredEnterpriseOrgs := make([]EnterpriseOrganization, 0)
	for _, enterpriseOrg := range enterpriseOrgs {
		exists := false
		for _, regionalOrganization := range regionalOrganizations {
			if regionalOrganization.Id == int64(enterpriseOrg.DatabaseId) {
				exists = true
				break
			}
		}
		if !exists {
			filteredEnterpriseOrgs = append(filteredEnterpriseOrgs, EnterpriseOrganization{Id: int64(enterpriseOrg.DatabaseId), Name: string(enterpriseOrg.Login)})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filteredEnterpriseOrgs)
}

func GetRegionalOrganizations(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	isEnabled := db.NullBool{Value: true}
	regionalOrganizations, err := db.SelectRegionalOrganization(&isEnabled)
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
	isEnabled := db.NullBool{Value: true}

	regionalOrganizations, total, err := db.SelectRegionalOrganizationByOption(offset, filter, search, orderBy, orderType, &isEnabled)
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
			IsRegionalOrganization:  regionalOrganization.IsRegionalOrganization,
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

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

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
		IsRegionalOrganization:  regionalOrganizationDto.IsRegionalOrganization,
		IsCleanUpMembersEnabled: regionalOrganizationDto.IsCleanUpMembersEnabled,
		IsIndexRepoEnabled:      regionalOrganizationDto.IsIndexRepoEnabled,
		IsCopilotRequestEnabled: regionalOrganizationDto.IsCopilotRequestEnabled,
		IsAccessRequestEnabled:  regionalOrganizationDto.IsAccessRequestEnabled,
		IsEnabled:               true,
		CreatedBy:               username,
	}

	_, err = db.InsertRegionalOrganization(regionalOrganization)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(regionalOrganizationDto)
}

func UpdateRegionalOrganization(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

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
		ModifiedBy:              username,
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

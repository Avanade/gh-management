package ghmgmt

import (
	"errors"
	"time"
)

type RegionalOrganization struct {
	Id                      int64
	Name                    string
	IsRegionalOrganization  bool
	IsIndexRepoEnabled      bool
	IsCopilotRequestEnabled bool
	IsAccessRequestEnabled  bool
	IsEnabled               bool
	CreatedBy               string
	Created                 time.Time
	ModifiedBy              string
	Modified                time.Time
}

type NullBool struct {
	Value bool
}

func SelectRegionalOrganization(isEnabled *NullBool) ([]RegionalOrganization, error) {
	db := ConnectDb()
	defer db.Close()

	param := make(map[string]interface{})
	if isEnabled != nil {
		param["IsEnabled"] = isEnabled.Value
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_RegionalOrganization_Select", param)
	if err != nil {
		return nil, err
	}

	regionalOrganizations := make([]RegionalOrganization, 0)
	for _, row := range result {
		organization := RegionalOrganization{
			Id:                      row["Id"].(int64),
			Name:                    row["Name"].(string),
			IsRegionalOrganization:  row["IsRegionalOrganization"].(bool),
			IsIndexRepoEnabled:      row["IsIndexRepoEnabled"].(bool),
			IsCopilotRequestEnabled: row["IsCopilotRequestEnabled"].(bool),
			IsAccessRequestEnabled:  row["IsAccessRequestEnabled"].(bool),
			IsEnabled:               row["IsEnabled"].(bool),
			Created:                 row["Created"].(time.Time),
		}
		if row["CreatedBy"] != nil {
			organization.CreatedBy = row["CreatedBy"].(string)
		}
		if row["ModifiedBy"] != nil {
			organization.ModifiedBy = row["ModifiedBy"].(string)
		}
		if row["Modified"] != nil {
			organization.Modified = row["Modified"].(time.Time)
		}
		regionalOrganizations = append(regionalOrganizations, organization)
	}

	return regionalOrganizations, nil
}

func SelectRegionalOrganizationByOption(offset, filter int64, search, orderBy, orderType string, isEnabled *NullBool) (regionalOrganizations []RegionalOrganization, total int64, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Offset":    offset,
		"Filter":    filter,
		"Search":    search,
		"OrderBy":   orderBy,
		"OrderType": orderType,
	}

	if isEnabled != nil {
		param["IsEnabled"] = isEnabled.Value
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_RegionalOrganization_Select_ByOption", param)
	if err != nil {
		return nil, 0, err
	}

	for _, row := range result {
		organization := RegionalOrganization{
			Id:                      row["Id"].(int64),
			Name:                    row["Name"].(string),
			IsRegionalOrganization:  row["IsRegionalOrganization"].(bool),
			IsIndexRepoEnabled:      row["IsIndexRepoEnabled"].(bool),
			IsCopilotRequestEnabled: row["IsCopilotRequestEnabled"].(bool),
			IsAccessRequestEnabled:  row["IsAccessRequestEnabled"].(bool),
			IsEnabled:               row["IsEnabled"].(bool),
			Created:                 row["Created"].(time.Time),
		}
		if row["CreatedBy"] != nil {
			organization.CreatedBy = row["CreatedBy"].(string)
		}
		if row["ModifiedBy"] != nil {
			organization.ModifiedBy = row["ModifiedBy"].(string)
		}
		if row["Modified"] != nil {
			organization.Modified = row["Modified"].(time.Time)
		}
		regionalOrganizations = append(regionalOrganizations, organization)
	}

	paramTotal := map[string]interface{}{
		"Search": search,
	}
	resultTotal, err := db.ExecuteStoredProcedureWithResult("usp_RegionalOrganization_Select_ByOption_Total", paramTotal)
	if err != nil {
		return nil, 0, err
	}

	total = resultTotal[0]["Total"].(int64)
	return regionalOrganizations, total, nil
}

func SelectRegionalOrganizationById(id int64) (*RegionalOrganization, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_RegionalOrganization_Select_ById", param)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("regional organization not found")
	}

	row := result[0]
	organization := RegionalOrganization{
		Id:                      row["Id"].(int64),
		Name:                    row["Name"].(string),
		IsRegionalOrganization:  row["IsRegionalOrganization"].(bool),
		IsIndexRepoEnabled:      row["IsIndexRepoEnabled"].(bool),
		IsCopilotRequestEnabled: row["IsCopilotRequestEnabled"].(bool),
		IsAccessRequestEnabled:  row["IsAccessRequestEnabled"].(bool),
		IsEnabled:               row["IsEnabled"].(bool),
		Created:                 row["Created"].(time.Time),
	}
	if row["CreatedBy"] != nil {
		organization.CreatedBy = row["CreatedBy"].(string)
	}
	if row["Modified"] != nil {
		organization.Modified = row["Modified"].(time.Time)
	}
	if row["ModifiedBy"] != nil {
		organization.ModifiedBy = row["ModifiedBy"].(string)
	}

	return &organization, nil
}

func UpdateRegionalOrganization(organization RegionalOrganization) error {
	db := ConnectDb()
	defer db.Close()

	err := ValidateRegionalOrganization(organization)
	if err != nil {
		return err
	}

	param := map[string]interface{}{
		"Id":                      organization.Id,
		"Name":                    organization.Name,
		"IsRegionalOrganization":  organization.IsRegionalOrganization,
		"IsIndexRepoEnabled":      organization.IsIndexRepoEnabled,
		"IsCopilotRequestEnabled": organization.IsCopilotRequestEnabled,
		"IsAccessRequestEnabled":  organization.IsAccessRequestEnabled,
		"IsEnabled":               organization.IsEnabled,
		"ModifiedBy":              organization.ModifiedBy,
	}

	_, err = db.ExecuteStoredProcedureWithResult("usp_RegionalOrganization_Update", param)
	if err != nil {
		return err
	}

	return nil
}

func InsertRegionalOrganization(organization RegionalOrganization) (int64, error) {
	db := ConnectDb()
	defer db.Close()

	err := ValidateRegionalOrganization(organization)
	if err != nil {
		return 0, err
	}

	param := map[string]interface{}{
		"Id":                      organization.Id,
		"Name":                    organization.Name,
		"IsRegionalOrganization":  organization.IsRegionalOrganization,
		"IsIndexRepoEnabled":      organization.IsIndexRepoEnabled,
		"IsCopilotRequestEnabled": organization.IsCopilotRequestEnabled,
		"IsAccessRequestEnabled":  organization.IsAccessRequestEnabled,
		"IsEnabled":               organization.IsEnabled,
		"CreatedBy":               organization.CreatedBy,
	}

	_, err = db.ExecuteStoredProcedureWithResult("usp_RegionalOrganization_Insert", param)
	if err != nil {
		return 0, err
	}

	return organization.Id, nil
}

func ValidateRegionalOrganization(organization RegionalOrganization) error {
	if organization.Id == 0 {
		return errors.New("id is required")
	}
	if organization.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

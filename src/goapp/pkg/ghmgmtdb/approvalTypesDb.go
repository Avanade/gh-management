package ghmgmt

import (
	"fmt"
	"strconv"
	"time"
)

type ApprovalType struct {
	Id                        int
	Name                      string
	ApproverUserPrincipalName string
	IsActive                  bool
	IsArchived                bool
	Created                   time.Time
	CreatedBy                 string
	Modified                  time.Time
	ModifiedBy                string
}

func GetAllActiveApprovers() interface{} {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_SelectAllActive", nil)
	if err != nil {
		return err
	}
	return result
}

func SelectApprovalTypes() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Select", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SelectApprovalTypesByFilter(offset, filter int, orderby, ordertype, search string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Offset":    offset,
		"Filter":    filter,
		"Search":    search,
		"OrderBy":   orderby,
		"OrderType": ordertype,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Select_ByFilter", param)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SelectTotalApprovalTypes() int {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_TotalCount", nil)
	total, err := strconv.Atoi(fmt.Sprint(result[0]["Total"]))
	if err != nil {
		return 0
	}
	return total
}

func SelectApprovalTypeById(id int) (*ApprovalType, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Select_ById", param)
	if err != nil {
		return nil, err
	}

	approvalType := ApprovalType{
		Id:                        int(result[0]["Id"].(int64)),
		Name:                      result[0]["Name"].(string),
		ApproverUserPrincipalName: result[0]["ApproverUserPrincipalName"].(string),
		IsActive:                  result[0]["IsActive"].(bool),
		IsArchived:                result[0]["IsArchived"].(bool),
	}

	if result[0]["Created"] != nil {
		approvalType.Created = result[0]["Created"].(time.Time)
	}

	if result[0]["CreatedBy"] != nil {
		approvalType.CreatedBy = result[0]["CreatedBy"].(string)
	}

	if result[0]["Modified"] != nil {
		approvalType.Modified = result[0]["Modified"].(time.Time)
	}

	if result[0]["ModifiedBy"] != nil {
		approvalType.ModifiedBy = result[0]["ModifiedBy"].(string)
	}

	return &approvalType, nil
}

func InsertApprovalType(approvalType ApprovalType) (int, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Name":                      approvalType.Name,
		"ApproverUserPrincipalName": approvalType.ApproverUserPrincipalName,
		"IsActive":                  approvalType.IsActive,
		"CreatedBy":                 approvalType.CreatedBy,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Insert", param)
	if err != nil {
		return 0, err
	}
	return int(result[0]["Id"].(int64)), nil
}

func UpdateApprovalType(approvalType ApprovalType) (int, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                        approvalType.Id,
		"Name":                      approvalType.Name,
		"ApproverUserPrincipalName": approvalType.ApproverUserPrincipalName,
		"IsActive":                  approvalType.IsActive,
		"ModifiedBy":                approvalType.CreatedBy,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Update_ById", param)
	if err != nil {
		return 0, err
	}
	return int(result[0]["Id"].(int64)), nil
}

func SetIsArchiveApprovalTypeById(approvalType ApprovalType) (int, bool, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                        approvalType.Id,
		"Name":                      approvalType.Name,
		"ApproverUserPrincipalName": approvalType.ApproverUserPrincipalName,
		"IsArchived":                approvalType.IsArchived,
		"ModifiedBy":                approvalType.ModifiedBy,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Update_IsArchived_ById", param)
	if err != nil {
		return 0, false, err
	}

	return int(result[0]["Id"].(int64)), result[0]["Status"].(bool), nil
}

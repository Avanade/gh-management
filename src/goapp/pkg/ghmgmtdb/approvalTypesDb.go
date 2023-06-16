package ghmgmt

import (
	"fmt"
	"strconv"
)

type ApprovalType struct {
	Id                        int
	Name                      string
	ApproverUserPrincipalName string
	IsActive                  bool
	IsArchived                bool
	Created                   string
	CreatedBy                 string
	Modified                  string
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

func SelectApprovalTypes() (interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Select", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SelectApprovalTypesByFilter(offset, filter int, orderby, ordertype, search string) (interface{}, error) {
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

func SelectApprovalTypeById(id int) (interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Select_ById", param)
	if err != nil {
		return nil, err
	}

	return &result[0], nil
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

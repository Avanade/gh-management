package ghmgmt

import (
	"time"
)

type ApprovalType struct {
	Id         int
	Name       string
	IsActive   bool
	IsArchived bool
	Created    time.Time
	CreatedBy  string
	Modified   time.Time
	ModifiedBy string
}

func GetAllActiveApprovers() interface{} {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryApprovalType_Select_Active", nil)
	if err != nil {
		return err
	}
	return result
}

func SelectApprovalTypeById(id int) (*ApprovalType, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryApprovalType_Select_ById", param)
	if err != nil {
		return nil, err
	}

	approvalType := ApprovalType{
		Id:         int(result[0]["Id"].(int64)),
		Name:       result[0]["Name"].(string),
		IsActive:   result[0]["IsActive"].(bool),
		IsArchived: result[0]["IsArchived"].(bool),
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
		"Name":      approvalType.Name,
		"IsActive":  approvalType.IsActive,
		"CreatedBy": approvalType.CreatedBy,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryApprovalType_Insert", param)
	if err != nil {
		return 0, err
	}
	return int(result[0]["Id"].(int64)), nil
}

func UpdateApprovalType(approvalType ApprovalType) (int, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":         approvalType.Id,
		"Name":       approvalType.Name,
		"IsActive":   approvalType.IsActive,
		"ModifiedBy": approvalType.CreatedBy,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryApprovalType_Update", param)
	if err != nil {
		return 0, err
	}
	return int(result[0]["Id"].(int64)), nil
}

func SetIsArchiveApprovalTypeById(approvalType ApprovalType) (int, bool, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":         approvalType.Id,
		"Name":       approvalType.Name,
		"IsArchived": approvalType.IsArchived,
		"ModifiedBy": approvalType.ModifiedBy,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryApprovalType_Update_IsArchived", param)
	if err != nil {
		return 0, false, err
	}

	return int(result[0]["Id"].(int64)), result[0]["Status"].(bool), nil
}

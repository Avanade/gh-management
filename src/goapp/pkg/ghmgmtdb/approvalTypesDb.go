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

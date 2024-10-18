package ghmgmt

import (
	"main/pkg/sql"
	"os"
)

func ConnectDb() *sql.DB {
	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	return db
}

func SearchCommunitiesProjectsUsers(searchText, offSet, filter, selectedSourceType, username string) ([]map[string]interface{}, int, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{
		"Search":            searchText,
		"OffSet":            offSet,
		"Filter":            filter,
		"UserPrincipalName": username,
	}

	if selectedSourceType != "" {
		params["SelectedSourceType"] = selectedSourceType
	}

	result, total, err := db.ExecuteStoredProcedureWithResultTotal("usp_Search", params)
	if err != nil {
		return nil, 0, err
	}

	return result, total, err
}

func UpdateProjectApprovalApproverResponse(itemId, remarks, responseDate, respondedBy string, approvalStatusId int) (bool, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{
		"ApprovalSystemGUID": itemId,
		"ApprovalStatusId":   approvalStatusId,
		"ApprovalRemarks":    remarks,
		"ApprovalDate":       responseDate,
		"RespondedBy":        respondedBy,
	}

	_, err := db.ExecuteStoredProcedure("usp_RepositoryApproval_Update_ApproverResponse", params)
	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateCommunityApprovalApproverResponse(itemId, remarks, responseDate string, approvalStatusId int) (bool, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{
		"ApprovalSystemGUID": itemId,
		"ApprovalStatusId":   approvalStatusId,
		"ApprovalRemarks":    remarks,
		"ApprovalDate":       responseDate,
	}

	_, err := db.ExecuteStoredProcedure("usp_ApprovalRequest_Update_CommunityApproverResponse", params)
	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateApprovalApproverResponse(itemId, remarks, responseDate string, approvalStatusId int, respondedBy string) (bool, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{
		"ApprovalSystemGUID": itemId,
		"ApprovalStatusId":   approvalStatusId,
		"ApprovalRemarks":    remarks,
		"ApprovalDate":       responseDate,
		"Approver":           respondedBy,
	}

	_, err := db.ExecuteStoredProcedure("usp_ApprovalRequest_Update_ApproverResponse", params)
	if err != nil {
		return false, err
	}

	return true, nil
}

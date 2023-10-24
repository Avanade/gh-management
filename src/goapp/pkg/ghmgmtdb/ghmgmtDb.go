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

func SearchCommunitiesProjectsUsers(searchText, offSet, rowCount, username string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{
		"searchText":    searchText,
		"offSet":        offSet,
		"rowCount":      rowCount,
		"userprincipal": username,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_Search_communities_projects_users", params)
	if err != nil {
		return nil, err
	}

	return result, err
}

func UpdateApprovalApproverResponse(storedProcedure, itemId, remarks, responseDate, respondedBy string, approvalStatusId int) (bool, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{
		"ApprovalSystemGUID": itemId,
		"ApprovalStatusId":   approvalStatusId,
		"ApprovalRemarks":    remarks,
		"ApprovalDate":       responseDate,
		"RespondedBy":        respondedBy,
	}

	_, err := db.ExecuteStoredProcedure(storedProcedure, params)
	if err != nil {
		return false, err
	}

	return true, nil
}
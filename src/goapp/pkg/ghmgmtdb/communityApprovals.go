package ghmgmt

import (
	"fmt"
)

func ApprovalInsert(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_CommunityApprovals_Insert", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func OrganizationApprovalInsert(params map[string]interface{}) error {
	db := ConnectDb()
	defer db.Close()

	_, err := db.ExecuteStoredProcedure("dbo.PR_OrganizationsApprovalRequests_Insert", params)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

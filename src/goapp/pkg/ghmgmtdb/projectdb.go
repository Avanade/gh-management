package ghmgmt

import (
	"fmt"
)

func ProjectApprovalsSelectById(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_ProjectApprovals_Select_ById", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}
func ProjectsSelectByUserPrincipalName(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_Select_ByUserPrincipalName", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}
func ProjectsApprovalUpdateApproverUserPrincipalName(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_ProjectsApproval_Update_ApproverUserPrincipalName", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunityApprovalslUpdateApproverUserPrincipalName(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_CommunityApprovals_Update_ApproverUserPrincipalName", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

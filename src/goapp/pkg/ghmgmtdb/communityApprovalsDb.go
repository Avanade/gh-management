package ghmgmt

func ApprovalInsert(approver string, description string, username string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{

		"ApproverUserPrincipalName": approver,
		"Name":                      description,
		"CreatedBy":                 username,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_CommunityApprovals_Insert", params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func OrganizationApprovalInsert(organizationId int, requestId int64) error {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{

		"OrganizationId": organizationId,
		"RequestId":      requestId,
	}

	_, err := db.ExecuteStoredProcedure("dbo.PR_OrganizationsApprovalRequests_Insert", params)
	if err != nil {
		return err
	}

	return nil
}

func CommunityApprovalInsert(communityId int, requestId int64) error {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{

		"CommunityId": communityId,
		"RequestId":   requestId,
	}

	_, err := db.ExecuteStoredProcedure("dbo.PR_CommunityApprovalRequests_Insert", params)
	if err != nil {
		return err
	}

	return nil
}

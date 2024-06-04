package ghmgmt

func ApprovalInsert(approver string, description string, username string) (requestId int64, err error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{

		"ApproverUserPrincipalName": approver,
		"Name":                      description,
		"CreatedBy":                 username,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_ApprovalRequest_Select_Insert", params)
	if err != nil {
		return
	}

	requestId = result[0]["Id"].(int64)
	return
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

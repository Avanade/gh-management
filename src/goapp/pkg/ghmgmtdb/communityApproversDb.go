package ghmgmt

func GetCommunityApprovers(category string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GuidanceCategory": category,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_CommunityApprover_Select", param)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetActiveCommunityApprovers(category string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GuidanceCategory": category,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_CommunityApprover_Select_Active", param)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateCommunityApproversById(id int, disabled bool, approverUserPrincipalName, username string, category string) (bool, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ApproverUserPrincipalName": approverUserPrincipalName,
		"IsDisabled":                disabled,
		"CreatedBy":                 username,
		"ModifiedBy":                username,
		"Id":                        id,
		"GuidanceCategory":          category,
	}

	_, err := db.ExecuteStoredProcedure("usp_CommunityApprover_Insert", param)
	if err != nil {
		return false, err
	}

	return true, nil
}

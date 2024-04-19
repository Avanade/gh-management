package ghmgmt

func GetCommunityApprovers(category string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Category": category,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_CommunityApproversList_select", param)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetActiveCommunityApprovers(category string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Category": category,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_CommunityApproversList_SelectAllActive", param)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetCommunityApproversById(id string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_CommunityApproversList_select_byId", param)
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
		"Disabled":                  disabled,
		"CreatedBy":                 username,
		"ModifiedBy":                username,
		"Id":                        id,
		"Category":                  category,
	}

	_, err := db.ExecuteStoredProcedure("PR_CommunityApproversList_Insert", param)
	if err != nil {
		return false, err
	}

	return true, nil
}

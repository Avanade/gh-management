package ghmgmt

func GetCommunityApprovers() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_CommunityApproversList_select", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetActiveCommunityApprovers() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_CommunityApproversList_SelectAllActive", nil)
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

func UpdateCommunityApproversById(id int, disabled bool, approverUserPrincipalName, username string) (bool, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ApproverUserPrincipalName": approverUserPrincipalName,
		"Disabled":                  disabled,
		"CreatedBy":                 username,
		"ModifiedBy":                username,
		"Id":                        id,
	}

	_, err := db.ExecuteStoredProcedure("PR_CommunityApproversList_Insert", param)
	if err != nil {
		return false, err
	}

	return true, nil
}

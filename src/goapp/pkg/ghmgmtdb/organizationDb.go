package ghmgmt

func OrganizationInsert(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Organizations_Insert", params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetRegionalOrganizationById(id int) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_RegionalOrganizations_SelectById", params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetAllRegionalOrganizations() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_RegionalOrganizations_SelectAll", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func OrganizationUpdateApprovalStatus(id int64, status int) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":               id,
		"ApprovalStatusID": status,
	}
	db.ExecuteStoredProcedure("PR_Organizations_UpdateStatus", param)
}

func GetAllOrganizationRequest(username string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Username": username,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Organizations_SelectByUser", param)
	if err != nil {
		return nil, err
	}

	return result, err
}

func GetOrganizationApprovalRequest(id int64) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_OrganizationsApprovalRequests_SelectByOrgId", param)
	if err != nil {
		return nil, err
	}

	return result, err
}

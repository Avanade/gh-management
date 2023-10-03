package ghmgmt

func SelectAllSponsors() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	sponsors, err := db.ExecuteStoredProcedureWithResult("PR_OSSContributionSponsors_SelectAll", nil)
	if err != nil {
		return nil, err
	}
	return sponsors, nil
}

func SelectSponsorsByIsArchived(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	sponsors, err := db.ExecuteStoredProcedureWithResult("PR_OSSContributionSponsors_SelectByIsArchived", params)
	if err != nil {
		return nil, err
	}
	return sponsors, nil
}

func InsertSponsor(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	sponsors, err := db.ExecuteStoredProcedureWithResult("PR_OSSContributionSponsors_Insert", params)
	if err != nil {
		return nil, err
	}
	return sponsors, nil

}

func UpdateSponsor(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	sponsors, err := db.ExecuteStoredProcedureWithResult("PR_OSSContributionSponsors_Update", params)
	if err != nil {
		return nil, err
	}
	return sponsors, nil
}

func SelectSponsorByName(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	sponsors, err := db.ExecuteStoredProcedureWithResult("PR_OSSContributionSponsors_SelectByName", params)
	if err != nil {
		return nil, err
	}
	return sponsors, nil
}

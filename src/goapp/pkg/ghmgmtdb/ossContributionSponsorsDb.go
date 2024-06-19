package ghmgmt

func SelectAllSponsors() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	sponsors, err := db.ExecuteStoredProcedureWithResult("usp_OSSContributionSponsor_Select", nil)
	if err != nil {
		return nil, err
	}
	return sponsors, nil
}

func SelectSponsorsByIsArchived(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	sponsors, err := db.ExecuteStoredProcedureWithResult("usp_OSSContributionSponsor_Select_ByIsArchived", params)
	if err != nil {
		return nil, err
	}
	return sponsors, nil
}

func InsertSponsor(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	sponsors, err := db.ExecuteStoredProcedureWithResult("usp_OSSContributionSponsor_Insert", params)
	if err != nil {
		return nil, err
	}
	return sponsors, nil

}

func UpdateSponsor(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	sponsors, err := db.ExecuteStoredProcedureWithResult("usp_OSSContributionSponsor_Update", params)
	if err != nil {
		return nil, err
	}
	return sponsors, nil
}

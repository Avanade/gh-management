package ghmgmt

func ExternalLinksExecuteSelect() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	ExternalLinks, err := db.ExecuteStoredProcedureWithResult("PR_ExternalLinks_Select", nil)
	if err != nil {
		return nil, err
	}
	return ExternalLinks, err
}

func ExternalLinksExecuteAllEnabled(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	ExternalLinks, err := db.ExecuteStoredProcedureWithResult("PR_ExternalLinks_SelectAllEnabled", params)
	if err != nil {
		return nil, err
	}
	return ExternalLinks, err
}

func ExternalLinksExecuteById(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	ExternalLinks, err := db.ExecuteStoredProcedureWithResult("PR_ExternalLinks_SelectById", params)
	if err != nil {
		return nil, err
	}
	return ExternalLinks, err
}

func ExternalLinksExecuteCreate(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	ExternalLinks, err := db.ExecuteStoredProcedureWithResult("PR_ExternalLinks_Insert", params)
	if err != nil {
		return nil, err
	}
	return ExternalLinks, err

}

func ExternalLinksExecuteUpdate(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	ExternalLinks, err := db.ExecuteStoredProcedureWithResult("PR_ExternalLinks_Update", params)
	if err != nil {
		return nil, err
	}
	return ExternalLinks, err

}

func ExternalLinksExecuteDelete(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	ExternalLinks, err := db.ExecuteStoredProcedureWithResult("PR_ExternalLinks_Delete", params)
	if err != nil {
		return nil, err
	}
	return ExternalLinks, err

}

package ghmgmt

func ExternalLinksExecuteSelect() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	externalLinks, err := db.ExecuteStoredProcedureWithResult("usp_ExternalLink_Select", nil)
	if err != nil {
		return nil, err
	}
	return externalLinks, err
}

func ExternalLinksExecuteAllEnabled(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	externalLinks, err := db.ExecuteStoredProcedureWithResult("usp_ExternalLink_Select_ByIsEnabled", params)
	if err != nil {
		return nil, err
	}
	return externalLinks, err
}

func ExternalLinksExecuteById(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	externalLinks, err := db.ExecuteStoredProcedureWithResult("usp_ExternalLink_Select_ById", params)
	if err != nil {
		return nil, err
	}
	return externalLinks, err
}

func ExternalLinksExecuteCreate(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	externalLinks, err := db.ExecuteStoredProcedureWithResult("usp_ExternalLink_Insert", params)
	if err != nil {
		return nil, err
	}
	return externalLinks, err

}

func ExternalLinksExecuteUpdate(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	externalLinks, err := db.ExecuteStoredProcedureWithResult("usp_ExternalLink_Update", params)
	if err != nil {
		return nil, err
	}
	return externalLinks, err

}

func ExternalLinksExecuteDelete(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	externalLinks, err := db.ExecuteStoredProcedureWithResult("usp_ExternalLink_Delete", params)
	if err != nil {
		return nil, err
	}
	return externalLinks, err

}

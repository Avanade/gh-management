package ghmgmt

func ADGroup_Insert(objectId string, ADGroup string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ObjectId": objectId,
		"ADGroup":  ADGroup,
	}

	_, err := db.ExecuteStoredProcedure("PR_GitHubAccess_Insert", param)
	if err != nil {
		return err
	}

	return nil
}

func ADGroup_SelectAll() ([]string, error) {
	db := ConnectDb()
	defer db.Close()

	var list []string
	result, err := db.ExecuteStoredProcedureWithResult("PR_GitHubAccess_SelectAll", nil)
	if err != nil {
		return nil, err
	}

	for _, group := range result {
		list = append(list, group["ObjectId"].(string))
	}

	return list, nil
}

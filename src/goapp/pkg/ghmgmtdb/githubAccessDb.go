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

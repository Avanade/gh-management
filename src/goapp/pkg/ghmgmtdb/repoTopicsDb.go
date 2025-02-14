package ghmgmt

type PopularTopics struct {
	Topic string
	Total int
}

func DeleteProjectTopics(projectId int) error {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"RepositoryId": projectId,
	}
	_, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryTopic_Delete", param)
	if err != nil {
		return err
	}
	return nil
}

func InsertProjectTopics(projectId int, topic string) error {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"RepositoryId": projectId,
		"Topic":        topic,
	}
	_, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryTopic_Insert", param)
	if err != nil {
		return err
	}
	return nil
}

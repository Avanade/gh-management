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

func GetPopularTopics(offset, rowCount int) ([]PopularTopics, error) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Offset":   offset,
		"RowCount": rowCount,
	}
	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryTopic_Select_PopularTopic", param)
	if err != nil {
		return nil, err
	}

	popularTopics := make([]PopularTopics, 0)

	for _, v := range result {
		popularTopic := PopularTopics{
			Topic: v["Topic"].(string),
			Total: int(v["Total"].(int64)),
		}
		popularTopics = append(popularTopics, popularTopic)
	}

	return popularTopics, nil
}

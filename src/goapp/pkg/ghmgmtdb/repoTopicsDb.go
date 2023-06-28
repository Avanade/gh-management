package ghmgmt

type PopularTopics struct {
	Topic string
	Total int
}

func DeleteProjectTopics(projectId int) error {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId": projectId,
	}
	_, err := db.ExecuteStoredProcedureWithResult("dbo.PR_RepoTopics_Delete_ByProjectId", param)
	if err != nil {
		return err
	}
	return nil
}

func InsertProjectTopics(projectId int, topic string) error {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId": projectId,
		"Topic":     topic,
	}
	_, err := db.ExecuteStoredProcedureWithResult("dbo.PR_RepoTopics_Insert", param)
	if err != nil {
		return err
	}
	return nil
}

func GetPopularTopics(offset, rowCount int) ([]PopularTopics, error) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"offset":   offset,
		"rowCount": rowCount,
	}
	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_RepoTopics_Select_PopularTopics", param)
	if err != nil {
		return nil, err
	}

	var popularTopics []PopularTopics

	for _, v := range result {
		popularTopic := PopularTopics{
			Topic: v["Topic"].(string),
			Total: int(v["Total"].(int64)),
		}
		popularTopics = append(popularTopics, popularTopic)
	}

	return popularTopics, nil
}

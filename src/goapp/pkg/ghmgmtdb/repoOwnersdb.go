package ghmgmt

import "main/models"

func RepoOwnersInsert(ProjectId int64, userPrincipalName string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId":         ProjectId,
		"UserPrincipalName": userPrincipalName,
	}

	_, err := db.ExecuteStoredProcedure("PR_RepoOwners_Insert", param)
	if err != nil {
		return err
	}

	return nil
}

func RepoOwnersByUserAndProjectId(id int64, userPrincipalName string) (RepoOwner []models.TypRepoOwner, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId":         id,
		"UserPrincipalName": userPrincipalName,
	}
	result, err := db.ExecuteStoredProcedureWithResult("PR_RepoOwners_Select_ByUserAndProjectId", param)
	if err != nil {
		println(err)
	}

	for _, v := range result {
		data := models.TypRepoOwner{
			Id:                v["ProjectId"].(int64),
			UserPrincipalName: v["UserPrincipalName"].(string),
		}
		RepoOwner = append(RepoOwner, data)
	}
	return RepoOwner, err
}

func SelectAllRepoNameAndOwners() (RepoOwner []models.TypRepoOwner, err error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_RepoOwners_SelectAllRepoNameAndOwners", nil)
	if err != nil {
		println(err)
	}

	for _, v := range result {
		data := models.TypRepoOwner{
			Id:                v["ProjectId"].(int64),
			RepoName:          v["Name"].(string),
			UserPrincipalName: v["UserPrincipalName"].(string),
		}
		RepoOwner = append(RepoOwner, data)
	}
	return RepoOwner, err
}

func GetRepoOwnersRecordByRepoId(id int64) (RepoOwner []models.TypRepoOwner, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId": id,
	}
	result, err := db.ExecuteStoredProcedureWithResult("PR_RepoOwners_Select_ByRepoId", param)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		data := models.TypRepoOwner{
			Id:                v["ProjectId"].(int64),
			UserPrincipalName: v["UserPrincipalName"].(string),
		}
		RepoOwner = append(RepoOwner, data)
	}
	return RepoOwner, nil
}

func GetRepoOwnersByProjectIdWithGHUsername(id int64) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_RepoOwners_SelectGHUser_ByRepoId", param)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func DeleteRepoOwnerRecordByUserAndProjectId(id int64, userPrincipalName string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId":         id,
		"UserPrincipalName": userPrincipalName,
	}
	_, err := db.ExecuteStoredProcedure("PR_RepoOwners_Delete_ByUserAndProjectId", param)
	if err != nil {
		return err
	}

	return nil
}

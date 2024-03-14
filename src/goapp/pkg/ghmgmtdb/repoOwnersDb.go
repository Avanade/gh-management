package ghmgmt

type RepoOwner struct {
	Id                int64
	RepoName          string
	UserPrincipalName string
}

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

func RepoOwnersByUserAndProjectId(id int64, userPrincipalName string) (repoOwner []RepoOwner, err error) {
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
		data := RepoOwner{
			Id:                v["ProjectId"].(int64),
			UserPrincipalName: v["UserPrincipalName"].(string),
		}
		repoOwner = append(repoOwner, data)
	}
	return repoOwner, err
}

func SelectAllRepoNameAndOwners() (repoOwner []RepoOwner, err error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_RepoOwners_SelectAllRepoNameAndOwners", nil)
	if err != nil {
		println(err)
	}

	for _, v := range result {
		data := RepoOwner{
			Id:                v["ProjectId"].(int64),
			RepoName:          v["Name"].(string),
			UserPrincipalName: v["UserPrincipalName"].(string),
		}
		repoOwner = append(repoOwner, data)
	}
	return repoOwner, err
}

func GetRepoOwnersRecordByRepoId(id int64) (repoOwner []RepoOwner, err error) {
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
		data := RepoOwner{
			Id:                v["ProjectId"].(int64),
			UserPrincipalName: v["UserPrincipalName"].(string),
		}
		repoOwner = append(repoOwner, data)
	}
	return repoOwner, nil
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

func CountOwnedRepoByVisibility(userPrincipalName, organization string, visibility int) (total int, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"UserPrincipalName": userPrincipalName,
		"Organization":      organization,
		"Visibility":        visibility,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_Projects_Count_OwnedRepoByVisibility", param)
	if err != nil {
		return 0, err
	}

	return int(result[0]["Total"].(int64)), nil
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

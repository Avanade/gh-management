package ghmgmt

import (
	"fmt"
	"strconv"
)

type RepoOwner struct {
	Id                int64
	RepoName          string
	UserPrincipalName string
}

func RepoOwnersInsert(ProjectId int64, userPrincipalName string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"RepositoryId":      ProjectId,
		"UserPrincipalName": userPrincipalName,
	}

	_, err := db.ExecuteStoredProcedure("usp_RepositoryOwner_Insert", param)
	if err != nil {
		return err
	}

	return nil
}

func IsOwner(id int64, userPrincipalName string) (isOwner bool, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"RepositoryId":      id,
		"UserPrincipalName": userPrincipalName,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryOwner_IsOwner", param)

	if err != nil {
		fmt.Println(err)
	}
	isExisting := strconv.FormatInt(result[0]["IsOwner"].(int64), 2)
	isOwner, _ = strconv.ParseBool(isExisting)

	return isOwner, err
}

func SelectAllRepoNameAndOwners() (repoOwner []RepoOwner, err error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryOwner_Select", nil)
	if err != nil {
		println(err)
	}

	for _, v := range result {
		data := RepoOwner{
			Id:                v["RepositoryId"].(int64),
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
		"@RepositoryId": id,
	}
	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryOwner_Select_ByRepositoryId", param)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		data := RepoOwner{
			Id:                v["RepositoryId"].(int64),
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
		"RepositoryId": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryOwner_Select_ByRepositoryId", param)
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

	result, err := db.ExecuteStoredProcedureWithResult("usp_Repository_TotalCountOwnedRepository_ByVisibility", param)
	if err != nil {
		return 0, err
	}

	return int(result[0]["Total"].(int64)), nil
}

func DeleteRepoOwnerRecordByUserAndProjectId(id int64, userPrincipalName string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"RepositoryId":      id,
		"UserPrincipalName": userPrincipalName,
	}
	_, err := db.ExecuteStoredProcedure("usp_RepositoryOwner_Delete", param)
	if err != nil {
		return err
	}

	return nil
}

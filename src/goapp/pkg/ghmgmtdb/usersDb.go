package ghmgmt

import "fmt"

func GetUsersWithGithub() interface{} {
	db := ConnectDb()
	defer db.Close()
	result, _ := db.ExecuteStoredProcedureWithResult("usp_User_Select_ByGitHubIdNotNull", nil)

	return result
}

func IsUserExist(userPrincipalName string) bool {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"UserPrincipalName": userPrincipalName,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("usp_User_IsExisting", param)

	return result[0]["Result"] == 1
}

func InsertUser(userPrincipalName, name, givenName, surName, jobTitle string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"UserPrincipalName": userPrincipalName,
		"Name":              name,
		"GivenName":         givenName,
		"SurName":           surName,
		"JobTitle":          jobTitle,
	}

	_, err := db.ExecuteStoredProcedure("usp_User_Insert", param)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserGithub(userPrincipalName, githubId, githubUser string, force int) (map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"UserPrincipalName": userPrincipalName,
		"GitHubId":          githubId,
		"GitHubUser":        githubUser,
		"Force":             force,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_User_Update", param)
	if err != nil {
		return nil, err
	}

	return result[0], nil
}

func Users_Get_GHUser(UserPrincipalName string) (GHUser string) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"UserPrincipalName": UserPrincipalName,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_User_Select_ByUserPrincipalName", param)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	if result == nil {
		return ""
	}

	if result[0]["GitHubUser"] == nil {
		return ""
	}
	GHUser = result[0]["GitHubUser"].(string)
	return GHUser
}

func GetUserByGitHubId(GitHubId string) ([]map[string]interface{}, error) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"GitHubId": GitHubId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_User_Select_ByGitHubId", param)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetUserByGitHubUsername(GitHubUser string) ([]map[string]interface{}, error) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"GitHubUser": GitHubUser,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_User_Select_ByGitHubUser", param)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetUserByUserPrincipal(UserPrincipalName string) ([]map[string]interface{}, error) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"UserPrincipalName": UserPrincipalName,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_User_Select_ByUserPrincipalName", param)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func UsersGetEmail(GithubUser string) (string, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GitHubUser": GithubUser,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_User_Select_ByGitHubUser", param)
	if err != nil {
		return "", err
	}
	if len(result) == 0 {
		return "", nil
	} else {
		return result[0]["UserPrincipalName"].(string), err
	}
}

func GetUserEmailByGithubId(GithubId string) (string, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GithubId": GithubId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_User_Select_ByGitHubId", param)
	if err != nil {
		return "", err
	}
	if len(result) == 0 {
		return "", nil
	} else {
		return result[0]["UserPrincipalName"].(string), err
	}
}

package ghmgmt

import "fmt"

func GetUsersWithGithub() interface{} {
	db := ConnectDb()
	defer db.Close()
	result, _ := db.ExecuteStoredProcedureWithResult("PR_Users_Select_WithGithub", nil)

	return result
}

func IsUserExist(userPrincipalName string) bool {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"UserPrincipalName": userPrincipalName,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Users_IsExisting", param)

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

	_, err := db.ExecuteStoredProcedure("PR_Users_Insert", param)
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

	result, err := db.ExecuteStoredProcedureWithResult("PR_Users_Update_GitHubUser", param)
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

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Users_Get_GHUser", param)

	if err != nil {
		fmt.Println(err)
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

	result, err := db.ExecuteStoredProcedureWithResult("PR_Users_Select_ByGitHubId", param)

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

	result, err := db.ExecuteStoredProcedureWithResult("PR_Users_Select_ByGitHubUsers", param)

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

	result, err := db.ExecuteStoredProcedureWithResult("PR_Users_Select_ByUserPrincipalName", param)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func IsUserAdmin(userPrincipalName string) bool {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"UserPrincipalName": userPrincipalName,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Admins_IsAdmin", param)

	return result[0]["Result"] == "1"
}

func UsersGetEmail(GithubUser string) (string, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GithubUser": GithubUser,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_Users_GetEmailByGitHubUsername", param)
	if err != nil {
		return "0", err
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

	result, err := db.ExecuteStoredProcedureWithResult("PR_Users_GetEmailByGitHubId", param)
	if err != nil {
		return "0", err
	}
	if len(result) == 0 {
		return "", nil
	} else {
		return result[0]["UserPrincipalName"].(string), err
	}
}

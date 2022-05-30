package ghmgmt

import (
	"fmt"
	"main/models"
	"main/pkg/sql"

	"os"
)

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

func ConnectDb() *sql.DB {
	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	return db
}

func PRProjectsInsert(body models.TypNewProjectReqBody, user string) {

	cp := sql.ConnectionParam{

		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	param := map[string]interface{}{

		"Name":                   body.Name,
		"CoOwner":                body.Coowner,
		"Description":            body.Description,
		"ConfirmAvaIP":           body.ConfirmAvaIP,
		"ConfirmEnabledSecurity": body.ConfirmSecIPScan,
		"CreatedBy":              user,
	}
	_, err := db.ExecuteStoredProcedure("dbo.PR_Projects_Insert", param)
	if err != nil {
		fmt.Println(err)
	}

}

func Projects_IsExisting(body models.TypNewProjectReqBody) bool {

	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	param := map[string]interface{}{

		"Name": body.Name,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_IsExisting", param)

	if err != nil {
		fmt.Println(err)
		return false
	}

	if result[0]["Result"] == "1" {
		return true
	} else {
		return false
	}

	return false
}

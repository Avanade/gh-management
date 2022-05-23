package ghmgmt

import (
	"fmt"
	models "main/models"
	sql "main/pkg/sql"
	"os"
)

func GetUsersWithGithub() interface{} {
	db := ConnectDb()
	defer db.Close()
	result, _ := db.ExecuteStoredProcedureWithResult("PR_Users_Select_WithGithub", nil)

	return result
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

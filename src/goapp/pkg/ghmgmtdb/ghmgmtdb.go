package ghmgmt

import (
	sql "main/pkg/sql"
	"os"
)

func GetUsersWithGithub() interface{} {
	db := ConnectDb()
	defer db.Close()
	result, _ := db.ExecuteStoredProcedureWithResult("PR_Users_Select_WithGithub", nil)

	return result
}

func ConnectDb() (*sql.DB) {
	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
		}

	db, _ := sql.Init(cp)

	return db
}
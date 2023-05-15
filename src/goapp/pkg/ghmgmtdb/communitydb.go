package ghmgmt

import (
	"database/sql"
	"fmt"
)

func CommunitiesSelectByID(params map[string]interface{}) (sql.Result, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedure("dbo.PR_Communities_select_byID", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

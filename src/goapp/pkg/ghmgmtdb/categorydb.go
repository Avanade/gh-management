package ghmgmt

import (
	"database/sql"
	"fmt"
)

func CategoryInsert(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Category_Insert", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CategoryArticlesInsert(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_CategoryArticles_Insert", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}
func CommunitiesSelectByID(params map[string]interface{}) (sql.Result, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedure("dbo.PR_Communities_select_byID", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

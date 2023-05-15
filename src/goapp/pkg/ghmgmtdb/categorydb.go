package ghmgmt

import (
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

func CategorySelect() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Category_select", nil)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}
func CategoryArticlesselectById(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_CategoryArticles_select_ById", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}
func CategoryArticlesSelectByArticlesID(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_CategoryArticles_select_ByArticlesID", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CategorySelectById(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Category_select_ById", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

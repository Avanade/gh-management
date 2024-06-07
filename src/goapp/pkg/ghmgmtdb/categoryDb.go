package ghmgmt

import (
	"fmt"
)

func CategoryInsert(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.usp_GuidanceCategory_Insert", params)
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

	result, err := db.ExecuteStoredProcedureWithResult("dbo.usp_GuidanceCategory_Select", nil)
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

	result, err := db.ExecuteStoredProcedureWithResult("dbo.usp_GuidanceCategory_Select_ById", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CategoryArticlesUpdate(params map[string]interface{}) error {
	db := ConnectDb()
	defer db.Close()

	_, err := db.ExecuteStoredProcedure("dbo.PR_CategoryArticles_Update", params)
	if err != nil {
		fmt.Println(err)
	}

	return err
}

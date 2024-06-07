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

	result, err := db.ExecuteStoredProcedureWithResult("usp_GuidanceCategoryArticle_Insert", params)
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

	result, err := db.ExecuteStoredProcedureWithResult("usp_GuidanceCategoryArticle_Select_ByCategoryId", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CategoryArticlesSelectByArticlesID(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_GuidanceCategoryArticle_Select_ByID", params)
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

	_, err := db.ExecuteStoredProcedure("usp_GuidanceCategoryArticle_Update", params)
	if err != nil {
		fmt.Println(err)
	}

	return err
}

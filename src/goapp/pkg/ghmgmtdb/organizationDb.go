package ghmgmt

import (
	"fmt"
)

func OrganizationInsert(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Organizations_Insert", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func GetRegionalOrganizationById(id int) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_RegionalOrganizations_SelectById", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

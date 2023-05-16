package ghmgmt

import (
	"database/sql"
	"fmt"
)

func CommunitiesSelectByID(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Communities_select_byID", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunitiesInsert(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Communities_Insert", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunitySponsorsInsert(params map[string]interface{}) (sql.Result, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedure("dbo.PR_CommunitySponsors_Insert", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func RelatedCommunitiesDelete(params map[string]interface{}) (sql.Result, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedure("dbo.PR_RelatedCommunities_Delete", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}
func RelatedCommunitiesInsert(params map[string]interface{}) (sql.Result, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedure("dbo.PR_RelatedCommunities_Insert", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunityTagsInsert(params map[string]interface{}) (sql.Result, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedure("dbo.PR_CommunityTags_Insert", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunitiesUpdate(params map[string]interface{}) (sql.Result, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedure("dbo.PR_Communities_Update", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}
func CommunityApprovalsSelectById(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_CommunityApprovals_Select_ById", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}
func CommunitiesIsexternal(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Communities_Isexternal", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}
func CommunitiesInitCommunityType(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Communities_InitCommunityType", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunitySponsorsSelect(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_CommunitySponsors_Select", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}
func CommunitySponsorsUpdate(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_CommunitySponsors_Update", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}
func CommunitySponsorsSelectByCommunityId(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_CommunitySponsors_Select_By_CommunityId", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}
func CommunityTagsSelectByCommunityId(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_CommunityTags_Select_By_CommunityId", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}
func RelatedCommunitiesSelect(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_RelatedCommunities_Select", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

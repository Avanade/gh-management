package ghmgmt

import (
	"fmt"
	"strconv"
)

type Activity struct {
	Name        string
	Url         string
	Date        string
	TypeId      int
	CommunityId int
	CreatedBy   string
	ModifiedBy  string
}

func CommunitiesActivities_Select() interface{} {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("usp_CommunityActivity_Select", nil)
	return result
}

func CommunitiesActivities_Select_ByOffsetAndFilterAndCreatedBy(offset, filter int, orderby, ordertype, search, createdBy string) interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Offset":    offset,
		"Filter":    filter,
		"Search":    search,
		"OrderType": ordertype,
		"OrderBy":   orderby,
		"CreatedBy": createdBy,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("usp_CommunityActivity_Select_ByOptionAndCreatedBy", param)
	return result
}

func CommunitiesActivities_Insert(body Activity) (int, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Name":           body.Name,
		"Url":            body.Url,
		"Date":           body.Date,
		"CreatedBy":      body.CreatedBy,
		"CommunityId":    body.CommunityId,
		"ActivityTypeId": body.TypeId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_CommunityActivity_Insert", param)
	if err != nil {
		return -1, err
	}
	id, err := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
	if err != nil {
		return -1, err
	}
	return id, nil
}

func CommunitiesActivities_TotalCount() int {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("usp_CommunityActivity_TotalCount_ByOption", nil)
	total, err := strconv.Atoi(fmt.Sprint(result[0]["Total"]))
	if err != nil {
		return 0
	}
	return total
}

func CommunitiesActivities_TotalCount_ByCreatedBy(createdBy, search string) int {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"CreatedBy": createdBy,
		"Search":    search,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("usp_CommunityActivity_TotalCount_ByOptionAndCreatedBy", param)
	total, err := strconv.Atoi(fmt.Sprint(result[0]["Total"]))
	if err != nil {
		return 0
	}
	return total
}

func CommunitiesActivities_Select_ById(id int) (interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_CommunityActivity_Select_ById", param)
	if err != nil {
		return nil, err
	}

	return &result[0], nil
}

func CommunityActivitiesHelpTypes_Insert(activityId int, helpTypeId int, details string) (int, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ActivityActivityId": activityId,
		"HelpTypeId":         helpTypeId,
		"Details":            details,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_CommunityActivityHelpType_Insert", param)
	if err != nil {
		return -1, err
	}

	id, err := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
	if err != nil {
		return -1, err
	}
	return id, nil
}

func ActivityTypes_Select() interface{} {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_ActivityType_Select", nil)
	if err != nil {
		return err
	}
	return result
}

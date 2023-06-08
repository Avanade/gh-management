package ghmgmt

import (
	"fmt"
	"main/models"
	"strconv"
)

func CommunityActivitiesContributionAreas_Insert(body models.CommunityActivitiesContributionAreas) (int, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"CommunityActivityId": body.CommunityActivityId,
		"ContributionAreaId":  body.ContributionAreaId,
		"IsPrimary":           body.IsPrimary,
		"CreatedBy":           body.CreatedBy,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_CommunityActivitiesContributionAreas_Insert", param)
	if err != nil {
		return -1, err
	}
	id, err := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
	if err != nil {
		return -1, err
	}
	return id, nil
}

func ContributionAreas_Select() (interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_ContributionAreas_Select", nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ContributionAreas_SelectByFilter(offset, filter int, orderby, ordertype, search string) (interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Offset":    offset,
		"Filter":    filter,
		"Search":    search,
		"OrderBy":   orderby,
		"OrderType": ordertype,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ContributionAreas_Select_ByFilter", param)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SelectTotalContributionAreas() int {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_ContributionAreas_TotalCount", nil)
	total, err := strconv.Atoi(fmt.Sprint(result[0]["Total"]))
	if err != nil {
		return 0
	}
	return total
}

func GetContributionAreaById(id int) interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ContributionAreas_SelectById", param)
	if err != nil {
		return err
	}
	return result
}

func UpdateContributionAreaById(id int, name string, username string) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":         id,
		"Name":       name,
		"ModifiedBy": username,
	}
	db.ExecuteStoredProcedure("PR_ContributionAreas_Update_ById", param)
}

func AdditionalContributionAreas_Select(activityId int) interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ActivityId": activityId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_AdditionalContributionAreas_Select_ByActivityId", param)
	if err != nil {
		return err
	}
	return result
}

func ContributionAreas_Insert(name, createdBy string) (int, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Name":      name,
		"CreatedBy": createdBy,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ContributionAreas_Insert", param)
	if err != nil {
		return -1, err
	}
	id, err := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
	if err != nil {
		return -1, err
	}
	return id, nil
}

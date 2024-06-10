package ghmgmt

import (
	"strconv"
	"strings"
	"time"
)

type OrganizationAccess struct {
	Id           int64
	User         User
	Organization RegionalOrganization
	Created      time.Time
}

type FailedCommunityApprovalRequestOrganizationAccess struct {
	Id                int64
	RegionName        string
	GitHubUsername    string
	UserPrincipalName string
	Approvers         []string
	RequestIds        []int64
}

type User struct {
	UserPrincipalName string
	GithubId          string
	GithubUsername    string
}

type RegionalOrganization struct {
	Id   int64
	Name string
}

func InsertOrganizationAccess(userPrincipalName string, organizationId int64) (id int64, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"UserPrincipalName": userPrincipalName,
		"OrganizationId":    organizationId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_OrganizationAccess_Insert", param)
	if err != nil {
		return
	}

	id = result[0]["Id"].(int64)

	return
}

func GetOrganizationAccessByUserPrincipalName(userPrincipalName string) ([]OrganizationAccess, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"UserPrincipalName": userPrincipalName,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_OrganizationAccess_SelectByUserPrincipalName", param)
	if err != nil {
		return nil, err
	}

	var organizationAccessArr []OrganizationAccess

	for _, v := range result {
		organizationAccess := OrganizationAccess{
			Id: v["Id"].(int64),
			User: User{
				UserPrincipalName: v["UserPrincipalName"].(string),
				GithubId:          v["GitHubId"].(string),
				GithubUsername:    v["GitHubUser"].(string),
			},
			Organization: RegionalOrganization{
				Id:   v["OrganizationId"].(int64),
				Name: v["OrganizationName"].(string),
			},
			Created: v["Created"].(time.Time),
		}

		organizationAccessArr = append(organizationAccessArr, organizationAccess)
	}

	return organizationAccessArr, err
}

func GetOrganizationAccessByApprovalRequestItemId(itemId string) (*OrganizationAccess, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ApprovalSystemGUID": itemId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_OrganizationAccess_SelectByApprovalRequestItemId", param)
	if err != nil {
		return nil, err
	}

	v := result[0]

	organizationAccess := OrganizationAccess{
		Id: v["Id"].(int64),
		User: User{
			UserPrincipalName: v["UserPrincipalName"].(string),
			GithubId:          v["GitHubId"].(string),
			GithubUsername:    v["GitHubUser"].(string),
		},
		Organization: RegionalOrganization{
			Id:   v["OrganizationId"].(int64),
			Name: v["OrganizationName"].(string),
		},
		Created: v["Created"].(time.Time),
	}

	return &organizationAccess, err
}

func HasOrganizationAccessPendingRequest(userPrincipalName string, organizationId int64) (hasPendingRequest bool, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"UserPrincipalName": userPrincipalName,
		"OrganizationId":    organizationId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_OrganizationAcess_HasPendingRequest", param)
	if err != nil {
		return
	}

	hasPendingRequest = result[0]["HasPendingRequest"].(bool)
	return
}

//  APPROVAL REQUESTS

func InsertOrganizationAccessApprovalRequest(organizationAccessId, requestId int64) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"OrganizationAccessId": organizationAccessId,
		"ApprovalRequestId":    requestId,
	}

	_, err := db.ExecuteStoredProcedure("usp_OrganizationAccessApprovalRequest_Insert", param)
	if err != nil {
		return err
	}

	return nil
}

func GetOrganizationAccessApprovalRequest(id int64) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_OrganizationAccessApprovalRequest_Select_ByOrganizationAccessId", param)
	if err != nil {
		return nil, err
	}

	return result, err
}

func GetFailedCommunityApprovalRequestOrganizationAccess() []FailedCommunityApprovalRequestOrganizationAccess {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("usp_ApprovalRequest_Select_FailedRequestOrganizationAccess", nil)

	var failedCommunityApprovalRequestOrganizationAccess []FailedCommunityApprovalRequestOrganizationAccess

	for _, v := range result {
		data := FailedCommunityApprovalRequestOrganizationAccess{
			Id:                v["Id"].(int64),
			RegionName:        v["RegionalOrgName"].(string),
			GitHubUsername:    v["GitHubUsername"].(string),
			UserPrincipalName: v["UserPrincipalName"].(string),
			Approvers:         []string{},
			RequestIds:        []int64{},
		}

		if v["Approvers"] != nil {
			approversStr := v["Approvers"].(string)
			data.Approvers = strings.Split(approversStr, ",")
		}

		if v["RequestIds"] != nil {
			requestIdsStr := v["RequestIds"].(string)
			requestIds := strings.Split(requestIdsStr, ",")

			for _, v := range requestIds {
				requestId, err := strconv.ParseInt(v, 0, 64)
				if err != nil {
					continue
				}
				data.RequestIds = append(data.RequestIds, int64(requestId))
			}
		}

		failedCommunityApprovalRequestOrganizationAccess = append(failedCommunityApprovalRequestOrganizationAccess, data)
	}

	return failedCommunityApprovalRequestOrganizationAccess
}

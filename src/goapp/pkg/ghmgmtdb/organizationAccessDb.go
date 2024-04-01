package ghmgmt

import "time"

type OrganizationAccess struct {
	Id           int64
	User         User
	Organization Organization
	Created      time.Time
}

type User struct {
	UserPrincipalName string
	GithubId          string
	GithubUsername    string
}

type Organization struct {
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
			Organization: Organization{
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
		Organization: Organization{
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
		"RequestId":            requestId,
	}

	_, err := db.ExecuteStoredProcedure("dbo.PR_OrganizationAccessApprovalRequest_Insert", param)
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

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_OrganizationAccessApprovalRequests_SelectByOrganizationAccessId", param)
	if err != nil {
		return nil, err
	}

	return result, err
}

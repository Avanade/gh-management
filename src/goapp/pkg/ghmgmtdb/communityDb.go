package ghmgmt

import (
	"database/sql"
	"fmt"
	"strconv"
)

type CommunityApproval struct {
	Id                         int64
	CommunityId                int64
	CommunityName              string
	CommunityUrl               string
	CommunityDescription       string
	CommunityNotes             string
	CommunityTradeAssocId      string
	CommunityType              string
	RequesterName              string
	RequesterGivenName         string
	RequesterSurName           string
	RequesterUserPrincipalName string
	ApproverUserPrincipalName  string
	ApprovalDescription        string
	ApproveUrl                 string
	RejectUrl                  string
	ApproveText                string
	RejectText                 string
}

type CommunityOnBoarding struct {
	Id                     int64              `json:"Id"`
	Name                   string             `json:"Name"`
	Url                    string             `json:"Url"`
	OnBoardingInstructions string             `json:"OnBoardingInstructions"`
	Sponsors               []CommunitySponsor `json:"Sponsors"`
	Communities            []RelatedCommunity `json:"Communities"`
}

type CommunitySponsor struct {
	Name      string `json:"Name"`
	GivenName string `json:"GivenName"`
	SurName   string `json:"SurName"`
	Email     string `json:"Email"`
}

type RelatedCommunity struct {
	Name       string `json:"Name"`
	Url        string `json:"Url"`
	IsExternal bool   `json:"IsExternal"`
}

func CommunitiesSelectByID(id string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_Community_Select_ById", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunitiesInsert(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_Community_Insert", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunitySponsorsInsert(params map[string]interface{}) (sql.Result, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedure("usp_CommunitySponsor_Insert", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func RelatedCommunitiesDelete(params map[string]interface{}) (sql.Result, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedure("usp_RelatedCommunity_Delete", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func RelatedCommunitiesInsert(params map[string]interface{}) (sql.Result, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedure("usp_RelatedCommunity_Insert", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunityTagsInsert(params map[string]interface{}) (sql.Result, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedure("usp_CommunityTag_Insert", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunityApprovalsSelectById(id int64) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_CommunityApprovals_Select_ById", param)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunitiesIsexternal(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_Community_IsExternal", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunitySponsorsSelectByCommunityId(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_CommunitySponsor_Select_ByCommunityId", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunityTagsSelectByCommunityId(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedureWithResult("usp_CommunityTag_Select_ByCommunityId", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func RelatedCommunitiesSelect(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()
	result, err := db.ExecuteStoredProcedureWithResult("usp_RelatedCommunity_Select_ByParentCommunityId", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func CommunityIManageExecuteSelect(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_Community_Select_IManage", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func MyCommunitites(username string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{
		"UserPrincipalName": username,
	}
	result, err := db.ExecuteStoredProcedureWithResult("usp_Community_Select_MyAllApproved", params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CommunityApprovalslUpdateApproverUserPrincipalName(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_ApprovalRequest_Update_ApproverUserPrincipalName", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func Communities_AddMember(CommunityId int, UserPrincipalName string) error {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"CommunityId":       CommunityId,
		"UserPrincipalName": UserPrincipalName,
	}

	_, err := db.ExecuteStoredProcedure("usp_CommunityMember_Insert", param)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func Communities_Related(CommunityId int64) (data []RelatedCommunity, err error) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"CommunityId": CommunityId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_Community_Select_RelatedCommunity", param)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range result {
		d := RelatedCommunity{
			Name:       v["Name"].(string),
			Url:        v["Url"].(string),
			IsExternal: v["IsExternal"].(bool),
		}
		data = append(data, d)
	}
	return
}

func Community_Sponsors(CommunityId int64) (data []CommunitySponsor, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"CommunityId": CommunityId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_CommunitySponsor_Select_ByCommunityId", param)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range result {
		d := CommunitySponsor{
			Name:      v["Name"].(string),
			GivenName: v["GivenName"].(string),
			SurName:   v["SurName"].(string),
			Email:     v["UserPrincipalName"].(string),
		}
		data = append(data, d)
	}
	return
}

func Community_Info(CommunityId int64) (data CommunityOnBoarding, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"Id": CommunityId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_Community_Select_ById", param)

	if err != nil {
		fmt.Println(err)
		return
	}

	data = CommunityOnBoarding{
		Id:                     result[0]["Id"].(int64),
		Name:                   result[0]["Name"].(string),
		OnBoardingInstructions: result[0]["OnBoardingInstructions"].(string),
		Url:                    result[0]["Url"].(string),
	}

	return
}

func Community_Onboarding_AddMember(CommunityId int64, UserPrincipalName string) (err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"CommunityId":       CommunityId,
		"UserPrincipalName": UserPrincipalName,
	}

	_, err = db.ExecuteStoredProcedure("usp_CommunityMember_Insert", param)

	if err != nil {
		fmt.Println(err)
	}
	return
}

func Community_Onboarding_RemoveMember(CommunityId int64, UserPrincipalName string) (err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"CommunityId":       CommunityId,
		"UserPrincipalName": UserPrincipalName,
	}

	_, err = db.ExecuteStoredProcedure("usp_CommunityMember_Delete", param)

	if err != nil {
		fmt.Println(err)
	}
	return
}

func Community_Membership_IsMember(CommunityId int64, UserPrincipalName string) (isMember bool, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"CommunityId":       CommunityId,
		"UserPrincipalName": UserPrincipalName,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_CommunityMember_IsExisting", param)

	if err != nil {
		fmt.Println(err)
	}
	isExisting := strconv.FormatInt(result[0]["IsExisting"].(int64), 2)
	isMember, _ = strconv.ParseBool(isExisting)
	return
}

func GetCommunities() interface{} {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_Community_Select_AllApproved", nil)
	if err != nil {
		return err
	}
	return result
}

func GetCommunityMembers(id int64) interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"CommunityId": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_CommunityMember_Select_ByCommunityId", param)
	if err != nil {
		return err
	}
	return result
}

func PopulateCommunityApproval(id int64) (CommunityApprovals []CommunityApproval) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"CommunityId": id,
	}
	result, _ := db.ExecuteStoredProcedureWithResult("usp_ApprovalRequest_Populate", param)

	for _, v := range result {
		data := CommunityApproval{
			Id:                         v["Id"].(int64),
			CommunityId:                v["CommunityId"].(int64),
			CommunityName:              v["CommunityName"].(string),
			CommunityUrl:               v["CommunityUrl"].(string),
			CommunityDescription:       v["CommunityDescription"].(string),
			CommunityNotes:             v["CommunityNotes"].(string),
			CommunityTradeAssocId:      v["CommunityTradeAssocId"].(string),
			CommunityType:              v["CommunityType"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			ApproverUserPrincipalName:  v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
		}
		CommunityApprovals = append(CommunityApprovals, data)
	}

	return
}

func GetFailedCommunityApprovalRequestCommunities() (CommunityApprovals []CommunityApproval) {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("usp_ApprovalRequest_Select_FailedRequestCommunity", nil)

	for _, v := range result {
		data := CommunityApproval{
			Id:                         v["Id"].(int64),
			CommunityId:                v["CommunityId"].(int64),
			CommunityName:              v["CommunityName"].(string),
			CommunityUrl:               v["CommunityUrl"].(string),
			CommunityDescription:       v["CommunityDescription"].(string),
			CommunityNotes:             v["CommunityNotes"].(string),
			CommunityTradeAssocId:      v["CommunityTradeAssocId"].(string),
			CommunityType:              v["CommunityType"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			ApproverUserPrincipalName:  v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
		}
		CommunityApprovals = append(CommunityApprovals, data)
	}

	return
}

func CommunityApprovalUpdateGUID(id int64, ApprovalSystemGUID string) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                 id,
		"ApprovalSystemGUID": ApprovalSystemGUID,
	}
	db.ExecuteStoredProcedure("usp_ApprovalRequest_Update_ApprovalSystemGUID", param)
}

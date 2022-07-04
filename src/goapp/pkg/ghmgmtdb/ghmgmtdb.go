package ghmgmt

import (
	"fmt"
	"main/models"
	"main/pkg/sql"
	"strconv"

	"os"
)

func GetUsersWithGithub() interface{} {
	db := ConnectDb()
	defer db.Close()
	result, _ := db.ExecuteStoredProcedureWithResult("PR_Users_Select_WithGithub", nil)

	return result
}

func IsUserExist(userPrincipalName string) bool {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"UserPrincipalName": userPrincipalName,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Users_IsExisting", param)

	return result[0]["Result"] == 1
}

func InsertUser(userPrincipalName, name, givenName, surName, jobTitle string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"UserPrincipalName": userPrincipalName,
		"Name":              name,
		"GivenName":         givenName,
		"SurName":           surName,
		"JobTitle":          jobTitle,
	}

	_, err := db.ExecuteStoredProcedure("PR_Users_Insert", param)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserGithub(userPrincipalName, githubId, githubUser string, force int) (map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"UserPrincipalName": userPrincipalName,
		"GitHubId":          githubId,
		"GitHubUser":        githubUser,
		"Force":             force,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_Users_Update_GitHubUser", param)
	if err != nil {
		return nil, err
	}

	return result[0], nil
}

func ConnectDb() *sql.DB {
	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	return db
}

// PROJECTS
func PRProjectsInsert(body models.TypNewProjectReqBody, user string) (id int64) {

	cp := sql.ConnectionParam{

		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)
	param := map[string]interface{}{

		"Name":                       body.Name,
		"CoOwner":                    body.Coowner,
		"Description":                body.Description,
		"ConfirmAvaIP":               body.ConfirmAvaIP,
		"ConfirmEnabledSecurity":     body.ConfirmSecIPScan,
		"CreatedBy":                  user,
		"Newcontribution":            body.Newcontribution,
		"OSSsponsor":                 body.OSSsponsor,
		"Avanadeofferingsassets":     body.Avanadeofferingsassets,
		"Willbecommercialversion":    body.Willbecommercialversion,
		"OSSContributionInformation": body.OSSContributionInformation,
	}
	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_Insert", param)
	if err != nil {
		fmt.Println(err)
	}
	id = result[0]["ItemId"].(int64)
	return
}

// PROJECTS
func PRProjectsUpdate(body models.TypNewProjectReqBody, user string) (id int64) {

	cp := sql.ConnectionParam{

		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}
	fmt.Println(body)
	db, _ := sql.Init(cp)
	param := map[string]interface{}{
		"ID":                         body.Id,
		"Name":                       body.Name,
		"CoOwner":                    body.Coowner,
		"Description":                body.Description,
		"ConfirmAvaIP":               body.ConfirmAvaIP,
		"ConfirmEnabledSecurity":     body.ConfirmSecIPScan,
		"ModifiedBy":                 user,
		"Newcontribution":            body.Newcontribution,
		"OSSsponsor":                 body.OSSsponsor,
		"Avanadeofferingsassets":     body.Avanadeofferingsassets,
		"Willbecommercialversion":    body.Willbecommercialversion,
		"OSSContributionInformation": body.OSSContributionInformation,
	}
	_, err := db.ExecuteStoredProcedure("dbo.PR_Projects_Update", param)
	if err != nil {
		fmt.Println(err)
	}

	return
}
func Projects_IsExisting(body models.TypNewProjectReqBody) bool {

	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	param := map[string]interface{}{

		"Name": body.Name,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_IsExisting", param)

	if err != nil {
		fmt.Println(err)
		return false
	}

	if result[0]["Result"] == "1" {
		return true
	} else {
		return false
	}
}

func PopulateProjectsApproval(id int64) (ProjectApprovals []models.TypProjectApprovals) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId": id,
	}
	result, _ := db.ExecuteStoredProcedureWithResult("PR_ProjectsApproval_Populate", param)
	for _, v := range result {
		data := models.TypProjectApprovals{
			Id:                         v["Id"].(int64),
			ProjectId:                  v["ProjectId"].(int64),
			ProjectName:                v["ProjectName"].(string),
			ProjectCoowner:             v["ProjectCoowner"].(string),
			ProjectDescription:         v["ProjectDescription"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			CoownerGivenName:           v["CoownerGivenName"].(string),
			CoownerSurName:             v["CoownerSurName"].(string),
			CoownerName:                v["CoownerName"].(string),
			CoownerUserPrincipalName:   v["CoownerUserPrincipalName"].(string),
			ApprovalTypeId:             v["ApprovalTypeId"].(int64),
			ApprovalType:               v["ApprovalType"].(string),
			ApproverUserPrincipalName:  v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
			Newcontribution:            v["newcontribution"].(string),
			OSSsponsor:                 v["OSSsponsor"].(string),
			Avanadeofferingsassets:     v["Avanadeofferingsassets"].(string),
			Willbecommercialversion:    v["Willbecommercialversion"].(string),
			OSSContributionInformation: v["OSSContributionInformation"].(string),
		}
		ProjectApprovals = append(ProjectApprovals, data)
	}

	return
}

func GetFailedProjectApprovalRequests() (ProjectApprovals []models.TypProjectApprovals) {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_Select_Failed", nil)

	for _, v := range result {
		data := models.TypProjectApprovals{
			Id:                         v["Id"].(int64),
			ProjectId:                  v["ProjectId"].(int64),
			ProjectName:                v["ProjectName"].(string),
			ProjectCoowner:             v["ProjectCoowner"].(string),
			ProjectDescription:         v["ProjectDescription"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			CoownerGivenName:           v["CoownerGivenName"].(string),
			CoownerSurName:             v["CoownerSurName"].(string),
			CoownerName:                v["CoownerName"].(string),
			CoownerUserPrincipalName:   v["CoownerUserPrincipalName"].(string),
			ApprovalTypeId:             v["ApprovalTypeId"].(int64),
			ApprovalType:               v["ApprovalType"].(string),
			ApproverUserPrincipalName:  v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
		}
		ProjectApprovals = append(ProjectApprovals, data)
	}

	return
}

func ProjectsApprovalUpdateGUID(id int64, ApprovalSystemGUID string) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                 id,
		"ApprovalSystemGUID": ApprovalSystemGUID,
	}
	db.ExecuteStoredProcedure("PR_ProjectsApproval_Update_ApprovalSystemGUID", param)
}

func GetProjectByName(projectName string) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Name": projectName,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Projects_Select_ByName", param)

	return result
}

func UpdateIsArchiveIsPrivate(projectName string, isArchived bool, isPrivate bool, username string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Name":       projectName,
		"IsArchived": isArchived,
		"IsPrivate":  isPrivate,
		"ModifiedBy": username,
	}

	_, err := db.ExecuteStoredProcedure("PR_Projects_Update_VisibilityByName", param)
	if err != nil {
		return err
	}

	return nil
}

// ACTIVITIES
func CommunitiesActivities_Select() interface{} {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_CommunityActivities_Select", nil)
	return result
}

func CommunitiesActivities_Select_ByOffsetAndFilter(offset, filter int, search string) interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Offset": offset,
		"Filter": filter,
		"Search": search,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_CommunityActivities_Select_ByOffsetAndFilter", param)
	return result
}

func CommunitiesActivities_Select_ByOffsetAndFilterAndCreatedBy(offset, filter int, search, createdBy string) interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Offset":    offset,
		"Filter":    filter,
		"Search":    search,
		"CreatedBy": createdBy,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_CommunityActivities_Select_ByOffsetAndFilterAndCreatedBy", param)
	return result
}

func CommunitiesActivities_Insert(body models.Activity) (int, error) {
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

	result, err := db.ExecuteStoredProcedureWithResult("PR_CommunityActivities_Insert", param)
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

	result, _ := db.ExecuteStoredProcedureWithResult("PR_CommunityActivities_TotalCount", nil)
	total, err := strconv.Atoi(fmt.Sprint(result[0]["Total"]))
	if err != nil {
		return 0
	}
	return total
}

func CommunitiesActivities_TotalCount_ByCreatedBy(createdBy string) int {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"CreatedBy": createdBy,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("[PR_CommunityActivities_TotalCount_ByCreatedBy]", param)
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

	result, err := db.ExecuteStoredProcedureWithResult("PR_CommunityActivities_Select_ById", param)
	if err != nil {
		return nil, err
	}

	return &result[0], nil
}

func ActivityTypes_Select() interface{} {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_ActivityTypes_Select", nil)
	if err != nil {
		return err
	}
	return result
}

func ActivityTypes_Insert(name string) (int, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Name": name,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ActivityTypes_Insert", param)
	if err != nil {
		return -1, err
	}
	id, err := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
	if err != nil {
		return -1, err
	}
	return id, nil
}

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

func ContributionAreas_Select() interface{} {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_ContributionAreas_Select", nil)
	if err != nil {
		return err
	}
	return result
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

//USERS
func Users_Get_GHUser(UserPrincipalName string) (GHUser string) {

	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	param := map[string]interface{}{

		"UserPrincipalName": UserPrincipalName,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Users_Get_GHUser", param)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	GHUser = result[0]["GitHubUser"].(string)
	return GHUser
}

func IsUserAdmin(userPrincipalName string) bool {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"UserPrincipalName": userPrincipalName,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Admins_IsAdmin", param)

	return result[0]["Result"] == "1"
}

// COMMUNITIES
func Communities_AddMember(CommunityId int, UserPrincipalName string) error {

	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	param := map[string]interface{}{
		"CommunityId":       CommunityId,
		"UserPrincipalName": UserPrincipalName,
	}

	_, err := db.ExecuteStoredProcedure("dbo.PR_CommunityMembers_Insert", param)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

func Communities_Related(CommunityId int64) (data []models.TypRelatedCommunities, err error) {

	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	param := map[string]interface{}{

		"CommunityId": CommunityId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Communities_Select_Related", param)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range result {
		d := models.TypRelatedCommunities{
			Name:       v["Name"].(string),
			Url:        v["Url"].(string),
			IsExternal: v["IsExternal"].(bool),
		}
		data = append(data, d)
	}
	return
}

func Community_Sponsors(CommunityId int64) (data []models.TypCommunitySponsorsList, err error) {
	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	param := map[string]interface{}{

		"CommunityId": CommunityId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_CommunitySponsors_Select_By_CommunityId", param)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range result {
		d := models.TypCommunitySponsorsList{
			Name:      v["Name"].(string),
			GivenName: v["GivenName"].(string),
			SurName:   v["SurName"].(string),
			Email:     v["UserPrincipalName"].(string),
		}
		data = append(data, d)
	}
	return
}

func Community_Info(CommunityId int64) (data models.TypCommunityOnBoarding, err error) {
	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	param := map[string]interface{}{

		"Id": CommunityId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Communities_select_byID", param)

	if err != nil {
		fmt.Println(err)
		return
	}

	data = models.TypCommunityOnBoarding{
		Id:   result[0]["Id"].(int64),
		Name: result[0]["Name"].(string),
		Url:  result[0]["Url"].(string),
	}

	return
}

func Community_Onboarding_AddMember(CommunityId int64, UserPrincipalName string) (err error) {
	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	param := map[string]interface{}{

		"CommunityId":       CommunityId,
		"UserPrincipalName": UserPrincipalName,
	}

	_, err = db.ExecuteStoredProcedure("dbo.PR_CommunityMembers_Insert", param)

	if err != nil {
		fmt.Println(err)
	}
	return
}

func Community_Onboarding_RemoveMember(CommunityId int64, UserPrincipalName string) (err error) {
	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	param := map[string]interface{}{

		"CommunityId":       CommunityId,
		"UserPrincipalName": UserPrincipalName,
	}

	_, err = db.ExecuteStoredProcedure("dbo.PR_CommunityMembers_Remove", param)

	if err != nil {
		fmt.Println(err)
	}
	return
}

func Community_Membership_IsMember(CommunityId int64, UserPrincipalName string) (isMember bool, err error) {
	cp := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, _ := sql.Init(cp)

	param := map[string]interface{}{

		"CommunityId":       CommunityId,
		"UserPrincipalName": UserPrincipalName,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_CommunityMembers_IsExisting", param)

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

	result, err := db.ExecuteStoredProcedureWithResult("PR_Communities_select", nil)
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

	result, err := db.ExecuteStoredProcedureWithResult("PR_CommunityMembers_Select_ByCommunityId", param)
	if err != nil {
		return err
	}
	return result
}
func PopulateCommunityApproval(id int64) (CommunityApprovals []models.TypCommunityApprovals) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"CommunityId": id,
	}
	result, _ := db.ExecuteStoredProcedureWithResult("PR_CommunityApprovals_Populate", param)

	for _, v := range result {
		data := models.TypCommunityApprovals{
			Id:                         v["Id"].(int64),
			CommunityId:                v["CommunityId"].(int64),
			CommunityName:              v["CommunityName"].(string),
			CommunityUrl:               v["CommunityUrl"].(string),
			CommunityDescription:       v["CommunityDescription"].(string),
			CommunityNotes:             v["CommunityNotes"].(string),
			CommunityTradeAssocId:      v["CommunityTradeAssocId"].(string),
			CommunityIsExternal:        v["CommunityIsExternal"].(bool),
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

func GetFailedCommunityApprovalRequests() (CommunityApprovals []models.TypCommunityApprovals) {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_CommunityApprovals_Select_Failed", nil)

	for _, v := range result {
		data := models.TypCommunityApprovals{
			Id:                         v["Id"].(int64),
			CommunityId:                v["CommunityId"].(int64),
			CommunityName:              v["CommunityName"].(string),
			CommunityUrl:               v["CommunityUrl"].(string),
			CommunityDescription:       v["CommunityDescription"].(string),
			CommunityNotes:             v["CommunityNotes"].(string),
			CommunityTradeAssocId:      v["CommunityTradeAssocId"].(string),
			CommunityIsExternal:        v["CommunityIsExternal"].(bool),
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
	db.ExecuteStoredProcedure("PR_CommunityApproval_Update_ApprovalSystemGUID", param)
}

// APPROVAL TYPES
func SelectApprovalTypes() (interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Select", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SelectApprovalTypesByFilter(offset, filter int, search string) (interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Offset": offset,
		"Filter": filter,
		"Search": search,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Select_ByFilter", param)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SelectTotalApprovalTypes() int {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_TotalCount", nil)
	total, err := strconv.Atoi(fmt.Sprint(result[0]["Total"]))
	if err != nil {
		return 0
	}
	return total
}

func SelectApprovalTypeById(id int) (interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Select_ById", param)
	if err != nil {
		return nil, err
	}

	return &result[0], nil
}

func InsertApprovalType(approvalType models.ApprovalType) (int, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Name":                      approvalType.Name,
		"ApproverUserPrincipalName": approvalType.ApproverUserPrincipalName,
		"IsActive":                  approvalType.IsActive,
		"CreatedBy":                 approvalType.CreatedBy,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Insert", param)
	if err != nil {
		return 0, err
	}
	return int(result[0]["Id"].(int64)), nil
}

func UpdateApprovalType(approvalType models.ApprovalType) (int, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                        approvalType.Id,
		"Name":                      approvalType.Name,
		"ApproverUserPrincipalName": approvalType.ApproverUserPrincipalName,
		"IsActive":                  approvalType.IsActive,
		"ModifiedBy":                approvalType.CreatedBy,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Update_ById", param)
	if err != nil {
		return 0, err
	}
	return int(result[0]["Id"].(int64)), nil
}

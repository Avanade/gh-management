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

		"Name":                   body.Name,
		"CoOwner":                body.Coowner,
		"Description":            body.Description,
		"ConfirmAvaIP":           body.ConfirmAvaIP,
		"ConfirmEnabledSecurity": body.ConfirmSecIPScan,
		"CreatedBy":              user,
	}
	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_Insert", param)
	if err != nil {
		fmt.Println(err)
	}
	id = result[0]["ItemId"].(int64)
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
func PRActivities_Select() interface{} {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Activities_Select", nil)
	return result
}

func PRActivities_Insert(name, url, createdBy string, communityId, activityId int) (int, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Name":           name,
		"Url":            url,
		"CreatedBy":      createdBy,
		"CommunityId":    communityId,
		"ActivityTypeId": activityId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_CommunitiesActivities_Insert", param)
	if err != nil {
		return -1, err
	}
	id, err := strconv.Atoi(fmt.Sprint(result[0]["Id"]))
	if err != nil {
		return -1, err
	}
	return id, nil
}

// ACTIVITIES TYPE
func PRActivityTypes_Select() interface{} {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_ActivityTypes_Select", nil)
	if err != nil {
		return err
	}
	return result
}

func PRActivityTypes_Insert(name string) (int, error) {
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

// CONTRIBUTION AREA
func PRContributionAreas_Select() interface{} {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_ContributionAreas_Select", nil)
	if err != nil {
		return err
	}
	return result
}

func PRContributionAreas_Insert(name, createdBy string) (int, error) {
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
			CommunityId:                  v["CommunityId"].(int64),
			CommunityName:                v["CommunityName"].(string),
			CommunityUrl:             v["CommunityUrl"].(string),
			CommunityDescription:         v["CommunityDescription"].(string),
			CommunityNotes:         v["CommunityNotes"].(string),
			CommunityTradeAssocId:         v["CommunityTradeAssocId"].(string),
			CommunityIsExternal:         v["CommunityIsExternal"].(bool),
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
			CommunityId:                  v["CommunityId"].(int64),
			CommunityName:                v["CommunityName"].(string),
			CommunityUrl:             v["CommunityUrl"].(string),
			CommunityDescription:         v["CommunityDescription"].(string),
			CommunityNotes:         v["CommunityNotes"].(string),
			CommunityTradeAssocId:         v["CommunityTradeAssocId"].(string),
			CommunityIsExternal:         v["CommunityIsExternal"].(bool),
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
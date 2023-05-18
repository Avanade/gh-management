package ghmgmt

import (
	"fmt"
	"main/models"
	"main/pkg/sql"
	"os"
	"strconv"
	"time"
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

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GithubId":                body.GithubId,
		"Name":                    body.Name,
		"CoOwner":                 body.Coowner,
		"Description":             body.Description,
		"ConfirmAvaIP":            body.ConfirmAvaIP,
		"ConfirmEnabledSecurity":  body.ConfirmSecIPScan,
		"ConfirmNotClientProject": body.ConfirmNotClientProject,
		"CreatedBy":               user,
		"VisibilityId":            body.Visibility,
		"TFSProjectReference":     body.TFSProjectReference,
	}
	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_Insert", param)
	if err != nil {
		fmt.Println(err)
	}
	id = result[0]["ItemId"].(int64)
	return
}

func ProjectInsertByImport(param map[string]interface{}) error {
	db := ConnectDb()
	defer db.Close()

	_, err := db.ExecuteStoredProcedure("dbo.PR_Projects_Insert", param)
	if err != nil {
		return err
	}
	return nil
}

func ProjectUpdateByImport(param map[string]interface{}) error {
	db := ConnectDb()
	defer db.Close()

	_, err := db.ExecuteStoredProcedure("dbo.PR_Projects_Update_Repo_Info", param)
	if err != nil {
		return err
	}
	return nil
}

func PRProjectsUpdate(body models.TypNewProjectReqBody, user string) (id int64) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ID":                      body.Id,
		"Name":                    body.Name,
		"CoOwner":                 body.Coowner,
		"Description":             body.Description,
		"ConfirmAvaIP":            body.ConfirmAvaIP,
		"ConfirmEnabledSecurity":  body.ConfirmSecIPScan,
		"ConfirmNotClientProject": body.ConfirmNotClientProject,
		"ModifiedBy":              user,
	}
	_, err := db.ExecuteStoredProcedure("dbo.PR_Projects_Update", param)
	if err != nil {
		fmt.Println(err)
	}

	return
}

func PRProjectsUpdateLegalQuestions(body models.TypeMakeProjectPublicReqBody, user string) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                         body.Id,
		"Newcontribution":            body.Newcontribution,
		"OSSsponsor":                 body.OSSsponsor,
		"Avanadeofferingsassets":     body.Avanadeofferingsassets,
		"Willbecommercialversion":    body.Willbecommercialversion,
		"OSSContributionInformation": body.OSSContributionInformation,
		"ModifiedBy":                 user,
	}
	_, err := db.ExecuteStoredProcedure("PR_Projects_Update_LegalQuestions", param)
	if err != nil {
		fmt.Println(err)
	}

	return
}

func Projects_IsExisting(body models.TypNewProjectReqBody) bool {

	db := ConnectDb()
	defer db.Close()

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

func Projects_IsExisting_By_GithubId(body models.TypNewProjectReqBody) bool {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GithubId": body.GithubId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_IsExisting_By_GithubId", param)

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
			ProjectDescription:         v["ProjectDescription"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			ApprovalTypeId:             v["ApprovalTypeId"].(int64),
			ApprovalType:               v["ApprovalType"].(string),
			ApproverUserPrincipalName:  v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
			Newcontribution:            v["newcontribution"].(string),
			OSSsponsor:                 v["OSSsponsor"].(string),
			Avanadeofferingsassets:     v["Avanadeofferingsassets"].(string),
			Willbecommercialversion:    v["Willbecommercialversion"].(string),
			OSSContributionInformation: v["OSSContributionInformation"].(string),
			RequestStatus:              v["RequestStatus"].(string),
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
			Newcontribution:            v["newcontribution"].(string),
			OSSsponsor:                 v["OSSsponsor"].(string),
			Avanadeofferingsassets:     v["Avanadeofferingsassets"].(string),
			Willbecommercialversion:    v["Willbecommercialversion"].(string),
			OSSContributionInformation: v["OSSContributionInformation"].(string),
			RequestStatus:              v["RequestStatus"].(string),
		}
		ProjectApprovals = append(ProjectApprovals, data)
	}

	return
}

func GetProjectApprovalsByProjectId(id int64) (ProjectApprovals []models.TypProjectApprovals) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_Select_By_ProjectId", param)

	for _, v := range result {
		data := models.TypProjectApprovals{
			Id:                        v["Id"].(int64),
			ProjectId:                 v["ProjectId"].(int64),
			ProjectName:               v["ProjectName"].(string),
			ApprovalTypeId:            v["ApprovalTypeId"].(int64),
			ApprovalType:              v["ApprovalType"].(string),
			ApproverUserPrincipalName: v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:       v["ApprovalDescription"].(string),
			RequestStatus:             v["RequestStatus"].(string),
		}
		ProjectApprovals = append(ProjectApprovals, data)
	}

	return
}

func GetProjectApprovalByGUID(id string) (ProjectApproval models.TypProjectApprovals) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ApprovalSystemGUID": id,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_Select_By_ApprovalSystemGUID", param)

	for _, v := range result {
		data := models.TypProjectApprovals{
			Id:                        v["Id"].(int64),
			ProjectId:                 v["ProjectId"].(int64),
			ProjectName:               v["ProjectName"].(string),
			ApprovalTypeId:            v["ApprovalTypeId"].(int64),
			ApprovalType:              v["ApprovalType"].(string),
			ApproverUserPrincipalName: v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:       v["ApprovalDescription"].(string),
			RequestStatus:             v["RequestStatus"].(string),
		}
		ProjectApproval = data
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

func GetProjectForRepoOwner() (RepoOwner []models.TypRepoOwner) {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Projects_ToRepoOwners", nil)

	for _, v := range result {
		data := models.TypRepoOwner{
			Id:                v["Id"].(int64),
			UserPrincipalName: v["UserPrincipalName"].(string),
		}
		RepoOwner = append(RepoOwner, data)
	}
	return RepoOwner
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

func GetProjectById(id int64) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Projects_Select_ById", param)

	return result
}

func GetProjectByGithubId(githubId int64) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GithubId": githubId,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Projects_Select_ByGithubId", param)

	return result
}

func Repos_Select_ByOffsetAndFilter(offset int, search string) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Offset": offset,
		"Search": search,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Repositories_Select_ByOffsetAndFilter", param)
	return result
}

func Repos_TotalCount_BySearchTerm(search string) int {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"search": search,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Repositories_TotalCount_BySearchTerm", param)
	total, err := strconv.Atoi(fmt.Sprint(result[0]["Total"]))
	if err != nil {
		return 0
	}
	return total
}

func UpdateProjectIsArchived(id int64, isArchived bool) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":         id,
		"IsArchived": isArchived,
	}

	_, err := db.ExecuteStoredProcedure("PR_Projects_Update_IsArchived_ById", param)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProjectVisibilityId(id int64, visibilityId int64) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":           id,
		"VisibilityId": visibilityId,
	}

	_, err := db.ExecuteStoredProcedure("PR_Projects_Update_Visibility_ById", param)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTFSProjectReferenceById(id int64, tFSProjectReference string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                  id,
		"TFSProjectReference": tFSProjectReference,
	}

	_, err := db.ExecuteStoredProcedure("PR_Projects_Update_TFSProjectReference_ById", param)
	if err != nil {
		return err
	}

	return nil
}

func GetRequestedReposByDateRange(start time.Time, end time.Time) ([]models.TypBasicRepo, error) {
	var RequestedRepos []models.TypBasicRepo
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Start": start.Format("2006-01-02"),
		"End":   end.Format("2006-01-02"),
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_Projects_Select_By_DateRange", param)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		data := models.TypBasicRepo{
			Name:        v["Name"].(string),
			Requestor:   v["CreatedBy"].(string),
			Description: v["Description"].(string),
		}
		RequestedRepos = append(RequestedRepos, data)
	}

	return RequestedRepos, nil
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

func CommunitiesActivities_TotalCount_ByCreatedBy(createdBy, search string) int {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"CreatedBy": createdBy,
		"Search":    search,
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

func CommunityActivitiesHelpTypes_Insert(activityId int, helpTypeId int, details string) (int, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ActivityActivityId": activityId,
		"HelpTypeId":         helpTypeId,
		"Details":            details,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_CommunityActivitiesHelpTypes_Insert", param)
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

func GetAllActiveApprovers() interface{} {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_SelectAllActive", nil)
	if err != nil {
		return err
	}
	return result
}

// USERS
func Users_Get_GHUser(UserPrincipalName string) (GHUser string) {

	db := ConnectDb()
	defer db.Close()

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

func GetUserByGitHubId(GitHubId string) ([]map[string]interface{}, error) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"GitHubId": GitHubId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_Users_Select_ByGitHubId", param)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetUserByGitHubUsername(GitHubUser string) ([]map[string]interface{}, error) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"GitHubUser": GitHubUser,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_Users_Select_ByGitHubUsers", param)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetUserByUserPrincipal(UserPrincipalName string) ([]map[string]interface{}, error) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"UserPrincipalName": UserPrincipalName,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_Users_Select_ByUserPrincipalName", param)

	if err != nil {
		return nil, err
	}

	return result, nil
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

	db := ConnectDb()
	defer db.Close()

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

	db := ConnectDb()
	defer db.Close()

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
	db := ConnectDb()
	defer db.Close()

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
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"Id": CommunityId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Communities_select_byID", param)

	if err != nil {
		fmt.Println(err)
		return
	}

	data = models.TypCommunityOnBoarding{
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

	_, err = db.ExecuteStoredProcedure("dbo.PR_CommunityMembers_Insert", param)

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

	_, err = db.ExecuteStoredProcedure("dbo.PR_CommunityMembers_Remove", param)

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

func SelectApprovalTypesByFilter(offset, filter int, orderby, ordertype, search string) (interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Offset":    offset,
		"Filter":    filter,
		"Search":    search,
		"OrderBy":   orderby,
		"OrderType": ordertype,
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

func SetIsArchiveApprovalTypeById(approvalType models.ApprovalType) (int, bool, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                        approvalType.Id,
		"Name":                      approvalType.Name,
		"ApproverUserPrincipalName": approvalType.ApproverUserPrincipalName,
		"IsArchived":                approvalType.IsArchived,
		"ModifiedBy":                approvalType.ModifiedBy,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalTypes_Update_IsArchived_ById", param)
	if err != nil {
		return 0, false, err
	}

	return int(result[0]["Id"].(int64)), result[0]["Status"].(bool), nil
}

func UsersGetEmail(GithubUser string) (string, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GithubUser": GithubUser,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_Users_GetEmailByGitHubUsername", param)
	if err != nil {
		return "0", err
	}
	if len(result) == 0 {
		return "", nil
	} else {
		return result[0]["UserPrincipalName"].(string), err
	}

}

func RepoOwnersInsert(ProjectId int64, userPrincipalName string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId":         ProjectId,
		"UserPrincipalName": userPrincipalName,
	}

	_, err := db.ExecuteStoredProcedure("PR_RepoOwners_Insert", param)
	if err != nil {
		return err
	}

	return nil
}

func RepoOwnersByUserAndProjectId(id int64, userPrincipalName string) (RepoOwner []models.TypRepoOwner, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId":         id,
		"UserPrincipalName": userPrincipalName,
	}
	result, err := db.ExecuteStoredProcedureWithResult("PR_RepoOwners_Select_ByUserAndProjectId", param)
	if err != nil {
		println(err)
	}

	for _, v := range result {
		data := models.TypRepoOwner{
			Id:                v["ProjectId"].(int64),
			UserPrincipalName: v["UserPrincipalName"].(string),
		}
		RepoOwner = append(RepoOwner, data)
	}
	return RepoOwner, err

}

func SelectAllRepoNameAndOwners() (RepoOwner []models.TypRepoOwner, err error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_RepoOwners_SelectAllRepoNameAndOwners", nil)
	if err != nil {
		println(err)
	}

	for _, v := range result {
		data := models.TypRepoOwner{
			Id:                v["ProjectId"].(int64),
			RepoName:          v["Name"].(string),
			UserPrincipalName: v["UserPrincipalName"].(string),
		}
		RepoOwner = append(RepoOwner, data)
	}
	return RepoOwner, err

}

func GetRepoOwnersRecordByRepoId(id int64) (RepoOwner []models.TypRepoOwner, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId": id,
	}
	result, err := db.ExecuteStoredProcedureWithResult("PR_RepoOwners_Select_ByRepoId", param)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		data := models.TypRepoOwner{
			Id:                v["ProjectId"].(int64),
			UserPrincipalName: v["UserPrincipalName"].(string),
		}
		RepoOwner = append(RepoOwner, data)
	}
	return RepoOwner, nil
}

func GetGitHubRepositories() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_Projects_SelectAllGitHub", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetRepoOwnersByProjectIdWithGHUsername(id int64) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_RepoOwners_SelectGHUser_ByRepoId", param)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func DeleteRepoOwnerRecordByUserAndProjectId(id int64, userPrincipalName string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId":         id,
		"UserPrincipalName": userPrincipalName,
	}
	_, err := db.ExecuteStoredProcedure("PR_RepoOwners_Delete_ByUserAndProjectId", param)
	if err != nil {
		return err
	}

	return nil
}

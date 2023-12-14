package main

import (
	"net/http"

	rtApi "main/routes/api"
	rtAzure "main/routes/login/azure"
	rtGithub "main/routes/login/github"
	rtPages "main/routes/pages"
	rtActivities "main/routes/pages/activity"
	rtAdmin "main/routes/pages/admin"
	rtCommunity "main/routes/pages/community"
	rtGuidance "main/routes/pages/guidance"
	rtProjects "main/routes/pages/project"
	rtSearch "main/routes/pages/search"

	"github.com/gorilla/mux"
)

func setPageRoutes(mux *mux.Router) {
	mux.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	mux.Handle("/", loadAzAuthPage(rtPages.HomeHandler))
	mux.Handle("/error/ghlogin", loadAzAuthPage(rtPages.GHLoginRequire))

	// SEARCH
	mux.Handle("/search", loadAzGHAuthPage(rtSearch.SearchHandler))

	// ACTIVITIES PAGE
	mux.Handle("/activities", loadAzGHAuthPage(rtActivities.IndexHandler))
	mux.Handle("/activities/{action:new}", loadAzGHAuthPage(rtActivities.FormHandler))
	mux.Handle("/activities/{action:edit|view}/{id}", loadAzGHAuthPage(rtActivities.FormHandler))

	// REPOSITORIES PAGE
	mux.Handle("/repositories", loadAzGHAuthPage(rtProjects.IndexHandler))
	mux.Handle("/repositories/new", loadAzGHAuthPage(rtProjects.FormHandler))
	mux.Handle("/repositories/makepublic/{id}", loadAzGHAuthPage(rtProjects.MakePublicHandler))

	// GUIDANCE PAGE
	mux.Handle("/guidance", loadAzGHAuthPage(rtGuidance.GuidanceHandler))
	mux.Handle("/guidance/categories/{id}", loadAzGHAuthPage(rtGuidance.EditCategoryHandler))
	mux.Handle("/guidance/articles/new", loadAzGHAuthPage(rtGuidance.NewArticleHandler))
	mux.Handle("/guidance/articles/{id}", loadAzGHAuthPage(rtGuidance.EditArticleHandler))

	// COMMUNITY PAGE
	mux.Handle("/communities", loadAzGHAuthPage(rtCommunity.CommunitiesHandler))
	mux.Handle("/communities/new", loadAzGHAuthPage(rtCommunity.CommunityFormHandler))
	mux.Handle("/communities/{id}", loadAzGHAuthPage(rtCommunity.CommunityFormHandler))
	mux.Handle("/communities/{id}/onboarding", loadAzGHAuthPage(rtCommunity.CommunityOnBoarding))

	// AUTHENTICATION
	mux.HandleFunc("/loginredirect", rtPages.LoginRedirectHandler)
	mux.HandleFunc("/gitredirect", rtPages.GitRedirectHandler)

	// AZURE
	mux.HandleFunc("/login/azure", rtAzure.LoginHandler)
	mux.HandleFunc("/login/azure/callback", rtAzure.CallbackHandler)
	mux.HandleFunc("/logout/azure", rtAzure.LogoutHandler)
	mux.HandleFunc("/authentication/azure/inprogress", rtPages.AuthenticationInProgressHandler)
	mux.HandleFunc("/authentication/azure/successful", rtPages.AuthenticationSuccessfulHandler)
	mux.HandleFunc("/authentication/azure/failed", rtPages.AuthenticationFailedHandler)

	// GITHUB
	mux.HandleFunc("/login/github", rtGithub.GithubLoginHandler)
	mux.HandleFunc("/login/github/callback", rtGithub.GithubCallbackHandler)
	mux.HandleFunc("/login/github/force", rtGithub.GithubForceSaveHandler)
	mux.HandleFunc("/logout/github", rtGithub.GitHubLogoutHandler)
	mux.HandleFunc("/authentication/github/inprogress", rtPages.GHAuthenticationInProgressHandler)
	mux.HandleFunc("/authentication/github/successful", rtPages.AuthenticationSuccessfulHandler)
	mux.HandleFunc("/authentication/github/failed", rtPages.AuthenticationFailedHandler)
}

func setAdminPageRoutes(mux *mux.Router) {
	// ADMIN
	muxAdmin := mux.PathPrefix("/admin").Subrouter()

	muxAdmin.Handle("", loadAdminPage(rtAdmin.AdminIndexHandler))
	muxAdmin.Handle("/members", loadAdminPage(rtAdmin.CommunityMembersHandler))
	muxAdmin.Handle("/guidance", loadAdminPage(rtGuidance.GuidanceHandler))
	muxAdmin.Handle("/communityapprovers", loadAdminPage(rtCommunity.CommunityApproversHandler))

	// APPROVAL TYPES ADMIN
	muxAdmin.Handle("/approvaltypes", loadAdminPage(rtAdmin.ApprovalTypesHandler))
	muxAdmin.Handle("/approvaltypes/{action:add}", loadAdminPage(rtAdmin.ApprovalTypeFormHandler))
	muxAdmin.Handle("/approvaltypes/{action:view|edit|delete}/{id}", loadAdminPage(rtAdmin.ApprovalTypeFormHandler))

	// CONTRIBUTION AREAS ADMIN
	muxAdmin.Handle("/contributionareas", loadAdminPage(rtAdmin.ContributionAreasHandler))
	muxAdmin.Handle("/contributionareas/{action:add}", loadAdminPage(rtAdmin.ContributionAreasFormHandler))
	muxAdmin.Handle("/contributionareas/{action:view|edit}/{id}", loadAdminPage(rtAdmin.ContributionAreasFormHandler))

	// EXTERNAL LINKS ADMIN
	muxAdmin.Handle("/externallinks", loadAdminPage(rtAdmin.ExternalLinksHandler))
	muxAdmin.Handle("/externallinks/{action:add}/", loadAdminPage(rtAdmin.ExternalLinksFormHandler))
	muxAdmin.Handle("/externallinks/{action:view|edit}/{id}", loadAdminPage(rtAdmin.ExternalLinksFormHandler))

	// OSS CONTRIBUTION SPONSORS ADMIN
	muxAdmin.Handle("/osscontributionsponsors", loadAdminPage(rtAdmin.OssContributionSponsorsHandler))
	muxAdmin.Handle("/osscontributionsponsors/form", loadAdminPage(rtAdmin.OssContributionSponsorsFormHandler))
}

func setApiRoutes(mux *mux.Router) {
	// APIS
	muxApi := mux.PathPrefix("/api").Subrouter()

	// ACTIVITIES API
	muxApi.Handle("/activity/type", loadAzGHAuthPage(rtApi.GetActivityTypes)).Methods("GET")
	muxApi.Handle("/activity/type", loadAzGHAuthPage(rtApi.CreateActivityType)).Methods("POST")
	muxApi.Handle("/activity", loadAzGHAuthPage(rtApi.CreateActivity)).Methods("POST")
	muxApi.Handle("/activity", loadAzGHAuthPage(rtApi.GetActivities)).Methods("GET")
	muxApi.Handle("/activity/{id}", loadAzGHAuthPage(rtApi.GetActivityById)).Methods("GET")

	// COMMUNITIES API
	muxApi.Handle("/community", loadAzGHAuthPage(rtApi.CommunityAPIHandler)).Methods("POST")
	muxApi.Handle("/communitySponsors", loadAzGHAuthPage(rtApi.CommunitySponsorsAPIHandler))
	muxApi.Handle("/CommunitySponsorsPerCommunityId/{id}", loadAzGHAuthPage(rtApi.CommunitySponsorsPerCommunityId))
	muxApi.Handle("/CommunityTagPerCommunityId/{id}", loadAzGHAuthPage(rtApi.CommunityTagPerCommunityId))
	muxApi.Handle("/community/onboarding/{id}", loadAzGHAuthPage(rtApi.GetCommunityOnBoardingInfo)).Methods("GET", "POST", "DELETE")
	muxApi.Handle("/community/all", loadAzAuthPage(rtApi.GetCommunities)).Methods("GET")
	muxApi.Handle("/community/{id}/members", loadAzAuthPage(rtApi.GetCommunityMembers)).Methods("GET")
	muxApi.Handle("/communitystatus/{id}", loadAzGHAuthPage(rtApi.GetRequestStatusByCommunity))
	muxApi.Handle("/community/getCommunitiesisexternal/{isexternal}", loadAzGHAuthPage(rtApi.GetCommunitiesIsexternal))
	muxApi.Handle("/community/members/processfile/{id}", loadAzGHAuthPage(rtApi.ProcessCommunityMembersListExcel)).Methods("POST")

	// RELATED COMMUNITIES API
	muxApi.Handle("/relatedcommunityAdd", loadAzAuthPage(rtApi.RelatedCommunitiesInsert))
	muxApi.Handle("/relatedcommunityDelete", loadAzAuthPage(rtApi.RelatedCommunitiesDelete))
	muxApi.Handle("/relatedcommunity/{id}", loadAzAuthPage(rtApi.RelatedCommunitiesSelect)).Methods("GET")

	// CONTRIBUTION AREAS API
	muxApi.Handle("/contributionarea", loadAzGHAuthPage(rtApi.CreateContributionAreas)).Methods("POST")
	muxApi.Handle("/contributionarea", loadAzGHAuthPage(rtApi.GetContributionAreas)).Methods("GET")
	muxApi.Handle("/contributionarea", loadAzGHAuthPage(rtApi.UpdateContributionArea)).Methods("PUT")
	muxApi.Handle("/contributionarea/{id}", loadAzGHAuthPage(rtApi.GetContributionAreaById)).Methods("GET")
	muxApi.Handle("/contributionarea/activity/{id}", loadAzGHAuthPage(rtApi.GetContributionAreasByActivityId)).Methods("GET")

	// CATEGORIES API
	muxApi.Handle("/Category", loadAzGHAuthPage(rtApi.CategoryAPIHandler)).Methods("POST")
	muxApi.Handle("/Category/list", loadAzGHAuthPage(rtApi.CategoryListAPIHandler))
	muxApi.Handle("/Category/update", loadAzGHAuthPage(rtApi.CategoryUpdate))
	muxApi.Handle("/Category/{id}", loadAzGHAuthPage(rtApi.GetCategoryByID))

	// CATEGORY ARTICLES API
	muxApi.Handle("/CategoryArticlesById/{id}", loadAzGHAuthPage(rtApi.GetCategoryArticlesById))
	muxApi.Handle("/CategoryArticlesByArticlesID/{id}", loadAzGHAuthPage(rtApi.GetCategoryArticlesByArticlesID))
	muxApi.Handle("/CategoryArticlesUpdate", loadAzGHAuthPage(rtApi.CategoryArticlesUpdate))

	// REPOSITORIES API
	muxApi.Handle("/repositories", loadAzGHAuthPage(rtApi.RequestRepository)).Methods("POST")
	muxApi.Handle("/repositories/{id}", loadAzGHAuthPage(rtApi.UpdateRepositoryById)).Methods("PUT")
	muxApi.Handle("/repositories/{id}/ecattid", loadAzGHAuthPage(rtApi.UpdateRepositoryEcattIdById)).Methods("PUT")
	muxApi.Handle("/repositories/list", loadAzGHAuthPage(rtApi.GetUserProjects))
	muxApi.Handle("/repositories/{id}", loadAzGHAuthPage(rtApi.GetRequestStatusByProject))
	muxApi.Handle("/repositories/request/public", loadAzGHAuthPage(rtApi.RequestMakePublic))
	muxApi.Handle("/repositories/collaborators/{id}", loadAzGHAuthPage(rtApi.GetRepoCollaboratorsByRepoId))
	muxApi.Handle("/repositories/collaborators/add/{id}/{ghUser}/{permission}", loadAzGHAuthPage(rtApi.AddCollaborator))
	muxApi.Handle("/repositories/collaborators/remove/{id}/{ghUser}/{permission}", loadAzGHAuthPage(rtApi.RemoveCollaborator))
	muxApi.Handle("/repositories/archive/{project}/{projectId}/{state}/{archive}", loadAzGHAuthPage(rtApi.ArchiveProject))
	muxApi.Handle("/repositories/visibility/{project}/{projectId}/{currentState}/{desiredState}", loadAzGHAuthPage(rtApi.SetVisibility))
	muxApi.Handle("/repositories/topics/popular", loadAzGHAuthPage(rtApi.GetPopularTopics))
	muxApi.Handle("/allusers", loadAzAuthPage(rtApi.GetAllUserFromActiveDirectory))
	muxApi.Handle("/alluserswithgithub", loadAzAuthPage(rtApi.GetUsersWithGithub))
	muxApi.Handle("/search/users/{search}", loadAzAuthPage(rtApi.SearchUserFromActiveDirectory))
	muxApi.Handle("/allrepositories", loadAzAuthPage(rtApi.GetAllRepositories))
	muxApi.Handle("/getActiveApprovalTypes", loadAzGHAuthPage(rtApi.GetActiveApprovalTypes))

	//APPROVAL TYPES API
	muxApi.Handle("/approval/type", loadAzAuthPage(rtApi.CreateApprovalType)).Methods("POST")
	muxApi.Handle("/approval/type/{id}", loadAzAuthPage(rtApi.EditApprovalTypeById)).Methods("PUT")
	muxApi.Handle("/approval/type/{id}/archived", loadAzAuthPage(rtApi.SetIsArchivedApprovalTypeById)).Methods("PUT")
	muxApi.Handle("/approval/types", loadAzAuthPage(rtApi.GetApprovalTypes)).Methods("GET")
	muxApi.Handle("/approval/type/{id}", loadAzAuthPage(rtApi.GetApprovalTypeById)).Methods("GET")

	//EXTERNAL LINKS API
	muxApi.Handle("/externallinks/create", loadAdminPage(rtApi.CreateExternalLinks)).Methods("POST")
	muxApi.Handle("/externallinks/update/{id}", loadAdminPage(rtApi.UpdateExternalLinks)).Methods("PUT")
	muxApi.Handle("/externallinks/", loadAdminPage(rtApi.GetExternalLinks))
	muxApi.Handle("/externallinks/enabled", loadAzAuthPage(rtApi.GetExternalLinksAllEnabled))
	muxApi.Handle("/externallinks/{id}", loadAdminPage(rtApi.GetExternalLinksById))
	muxApi.Handle("/externallinks/{action:delete}/{id}", loadAdminPage(rtApi.ExternalLinksDelete))

	// APPROVALS API
	muxApi.HandleFunc("/approvals/project/callback", rtApi.UpdateApprovalStatusProjects).Methods("POST")
	muxApi.HandleFunc("/approvals/project/reassign/callback", rtApi.UpdateApprovalReassignApprover)
	muxApi.Handle("/users/{username}/approvals", loadAzAuthPage(rtApi.DownloadProjectApprovalsByUsername))

	muxApi.HandleFunc("/approvals/community/reassign/callback", rtApi.UpdateCommunityApprovalReassignApprover)
	muxApi.HandleFunc("/approvals/community/callback", rtApi.UpdateApprovalStatusCommunity).Methods("POST")

	// COMMUNITY APPROVERS API
	muxApi.HandleFunc("/communityapprovers/update", rtCommunity.CommunityApproversListUpdate)
	muxApi.HandleFunc("/communityapprovers/GetCommunityApproversList", rtCommunity.GetCommunityApproversList)
	muxApi.HandleFunc("/communityapprovers/GetAllActiveCommunityApprovers", rtCommunity.GetAllActiveCommunityApprovers)
	muxApi.HandleFunc("/communityapprovers/GetCommunityApproversList/{id}", rtCommunity.GetCommunityApproversById)

	// OSS CONTRIBUTION SPONSORS API
	muxApi.Handle("/osscontributionsponsors/all", loadAdminPage((rtApi.GetAllOssContributionSponsors)))
	muxApi.Handle("/osscontributionsponsors/enabled", loadAzAuthPage((rtApi.GetAllEnabledOssContributionSponsors)))
	muxApi.Handle("/osscontributionsponsors/add", loadAdminPage((rtApi.AddSponsor)))
	muxApi.Handle("/osscontributionsponsors/update", loadAdminPage((rtApi.UpdateSponsor)))

	// API FOR LOGIC APP
	muxApi.Handle("/indexorgrepos", loadGuidAuthApi(rtApi.IndexOrgRepos)).Methods("GET")
	muxApi.Handle("/clearorgrepos", loadGuidAuthApi(rtApi.ClearOrgRepos)).Methods("GET")
	muxApi.Handle("/checkAvaInnerSource", loadGuidAuthApi(rtApi.CheckAvaInnerSource)).Methods("GET")
	muxApi.Handle("/checkAvaOpenSource", loadGuidAuthApi(rtApi.CheckAvaOpenSource)).Methods("GET")
	muxApi.Handle("/clearOrgMembers", loadGuidAuthApi(rtApi.ClearOrgMembers)).Methods("GET")
	muxApi.Handle("/RepoOwnerScan", loadGuidAuthApi(rtApi.RepoOwnerScan)).Methods("GET")
	muxApi.Handle("/RepoOwnersCleanup", loadGuidAuthApi(rtApi.RepoOwnersCleanup)).Methods("GET")
	muxApi.Handle("/recurringapproval", loadGuidAuthApi(rtApi.RecurringApproval)).Methods("GET")

	muxApi.Handle("/utility/fillout/approvers", loadGuidAuthApi(rtApi.FillOutApprovers)).Methods("GET")
	muxApi.Handle("/utility/fillout/approvalrequest/approvers", loadGuidAuthApi(rtApi.FillOutApprovalRequestApprovers)).Methods("GET")
	muxApi.Handle("/migrateOssSponsors", loadGuidAuthApi(rtApi.MigrateToOssSponsorsTable)).Methods("GET")

	// LEGACY APIS
	muxApi.Handle("/searchresult/{searchText}", loadGuidAuthApi(rtApi.LegacySearchHandler))
	mux.HandleFunc("/Home/Asset/{assetCode}", rtApi.RedirectAsset)
	mux.HandleFunc("/Home/AssetRequestCreation", rtApi.RedirectAssetRequest)
	mux.HandleFunc("/Home/AssetRequestCreation/", rtApi.RedirectAssetRequest)
	mux.Handle("/Home/Tool/{assetCode}", loadAzAuthPage(rtPages.ToolHandler))

	mux.Handle("/search/{offSet}/{rowCount}", loadAzGHAuthPage(rtSearch.GetSearchResults))
	mux.Handle("/communities/my", loadAzGHAuthPage(rtCommunity.GetMyCommunitylist))
	mux.Handle("/communities/imanage", loadAzGHAuthPage(rtCommunity.GetCommunityIManagelist))
	mux.Handle("/community/getcommunity/{id}", loadAzGHAuthPage(rtCommunity.GetUserCommunity))
	mux.Handle("/community", loadAzGHAuthPage(rtCommunity.GetUserCommunitylist))
}

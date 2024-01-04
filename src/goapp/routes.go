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
	mux.Handle("/guidance", loadAzGHAuthPage(rtGuidance.IndexHandler))
	mux.Handle("/guidance/categories/{id}", loadAzGHAuthPage(rtGuidance.EditCategoryHandler))
	mux.Handle("/guidance/articles/new", loadAzGHAuthPage(rtGuidance.NewArticleHandler))
	mux.Handle("/guidance/articles/{id}", loadAzGHAuthPage(rtGuidance.EditArticleHandler))

	// COMMUNITY PAGE
	mux.Handle("/communities", loadAzGHAuthPage(rtCommunity.IndexHandler))
	mux.Handle("/communities/new", loadAzGHAuthPage(rtCommunity.FormHandler))
	mux.Handle("/communities/{id}", loadAzGHAuthPage(rtCommunity.FormHandler))
	mux.Handle("/communities/{id}/onboarding", loadAzGHAuthPage(rtCommunity.OnBoardingHandler))

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

	// LEGACY REDIRECTS
	mux.HandleFunc("/Home/Asset/{assetCode}", rtApi.RedirectAsset)
	mux.HandleFunc("/Home/AssetRequestCreation", rtApi.RedirectAssetRequest)
	mux.HandleFunc("/Home/AssetRequestCreation/", rtApi.RedirectAssetRequest)
	mux.Handle("/Home/Tool/{assetCode}", loadAzAuthPage(rtPages.ToolHandler))
	mux.Handle("/search/{offSet}/{rowCount}", loadAzGHAuthPage(rtSearch.GetSearchResults))
}

func setAdminPageRoutes(mux *mux.Router) {
	// ADMIN
	muxAdmin := mux.PathPrefix("/admin").Subrouter()

	muxAdmin.Handle("", loadAdminPage(rtAdmin.AdminIndexHandler))

	// COMMUNITY MEMBERS
	muxAdmin.Handle("/members", loadAdminPage(rtAdmin.CommunityMembersHandler))

	// COMMUNITY APPROVERS
	muxAdmin.Handle("/communityapprovers", loadAdminPage(rtAdmin.CommunityApproversHandler))

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

	// ACTIVITY TYPES API
	muxApi.Handle("/activity-types", loadAzGHAuthPage(rtApi.GetActivityTypes)).Methods("GET")

	// ACTIVITY API
	muxApi.Handle("/activities", loadAzGHAuthPage(rtApi.CreateActivity)).Methods("POST")
	muxApi.Handle("/activities", loadAzGHAuthPage(rtApi.GetActivities)).Methods("GET")
	muxApi.Handle("/activities/{id}", loadAzGHAuthPage(rtApi.GetActivityById)).Methods("GET")

	// COMMUNITIES API
	muxApi.Handle("/communities", loadAzAuthPage(rtApi.GetCommunities)).Methods("GET")
	muxApi.Handle("/community/{id}/members", loadAzAuthPage(rtApi.GetCommunityMembers)).Methods("GET")
	muxApi.Handle("/communities/my", loadAzGHAuthPage(rtApi.GetMyCommunitylist)).Methods("GET")
	muxApi.Handle("/communities/imanage", loadAzGHAuthPage(rtApi.GetCommunityIManagelist)).Methods("GET")
	muxApi.Handle("/community/getcommunity/{id}", loadAzGHAuthPage(rtApi.GetUserCommunity)).Methods("GET")
	muxApi.Handle("/community/getCommunitiesisexternal/{isexternal}", loadAzGHAuthPage(rtApi.GetCommunitiesIsexternal)).Methods("GET")
	muxApi.Handle("/community", loadAzGHAuthPage(rtApi.AddCommunity)).Methods("POST")
	muxApi.Handle("/community/members/processfile/{id}", loadAzGHAuthPage(rtApi.ProcessCommunityMembersListExcel)).Methods("POST")
	muxApi.Handle("/communitySponsors", loadAzGHAuthPage(rtApi.CommunitySponsorsAPIHandler))
	muxApi.Handle("/CommunitySponsorsPerCommunityId/{id}", loadAzGHAuthPage(rtApi.CommunitySponsorsPerCommunityId))
	muxApi.Handle("/CommunityTagPerCommunityId/{id}", loadAzGHAuthPage(rtApi.CommunityTagPerCommunityId))
	muxApi.Handle("/community/onboarding/{id}", loadAzGHAuthPage(rtApi.GetCommunityOnBoardingInfo)).Methods("GET", "POST", "DELETE")
	muxApi.Handle("/communitystatus/{id}", loadAzGHAuthPage(rtApi.GetRequestStatusByCommunity))

	// RELATED COMMUNITIES API
	muxApi.Handle("/relatedcommunity/{id}", loadAzAuthPage(rtApi.RelatedCommunitiesSelect)).Methods("GET")

	muxApi.HandleFunc("/approvals/community/reassign/callback", rtApi.UpdateCommunityApprovalReassignApprover).Methods("POST")
	muxApi.HandleFunc("/approvals/community/callback", rtApi.UpdateApprovalStatusCommunity).Methods("POST")

	// COMMUNITY APPROVERS API
	muxApi.HandleFunc("/communityapprovers/update", rtApi.CommunityApproversListUpdate)
	muxApi.HandleFunc("/communityapprovers/GetCommunityApproversList", rtApi.GetCommunityApproversList)
	muxApi.HandleFunc("/communityapprovers/GetAllActiveCommunityApprovers", rtApi.GetAllActiveCommunityApprovers)
	muxApi.HandleFunc("/communityapprovers/GetCommunityApproversList/{id}", rtApi.GetCommunityApproversById)

	// CONTRIBUTION AREAS API
	muxApi.Handle("/contribution-areas", loadAzGHAuthPage(rtApi.CreateContributionAreas)).Methods("POST")
	muxApi.Handle("/contribution-areas", loadAzGHAuthPage(rtApi.GetContributionAreas)).Methods("GET")
	muxApi.Handle("/contribution-areas", loadAzGHAuthPage(rtApi.UpdateContributionArea)).Methods("PUT")
	muxApi.Handle("/contribution-areas/{id}", loadAzGHAuthPage(rtApi.GetContributionAreaById)).Methods("GET")
	muxApi.Handle("/activities/{id}/contribution-areas", loadAzGHAuthPage(rtApi.GetContributionAreasByActivityId)).Methods("GET")

	// CATEGORIES API
	muxApi.Handle("/categories", loadAzGHAuthPage(rtApi.CategoryAPIHandler)).Methods("POST")
	muxApi.Handle("/categories", loadAzGHAuthPage(rtApi.CategoryListAPIHandler)).Methods("GET")
	muxApi.Handle("/categories/{id}", loadAzGHAuthPage(rtApi.CategoryUpdate)).Methods("PUT")
	muxApi.Handle("/categories/{id}", loadAzGHAuthPage(rtApi.GetCategoryByID)).Methods("GET")

	// CATEGORY ARTICLES API
	muxApi.Handle("/categories/{id}/articles", loadAzGHAuthPage(rtApi.GetCategoryArticlesByCategoryId)).Methods("GET")
	muxApi.Handle("/articles/{id}", loadAzGHAuthPage(rtApi.GetCategoryArticlesById)).Methods("GET")
	muxApi.Handle("/articles/{id}", loadAzGHAuthPage(rtApi.UpdateCategoryArticlesById)).Methods("PUT")

	// REPOSITORIES API
	muxApi.Handle("/allrepositories", loadAzAuthPage(rtApi.GetAllRepositories))
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
	muxApi.Handle("/getActiveApprovalTypes", loadAzGHAuthPage(rtApi.GetActiveApprovalTypes))

	//APPROVAL TYPES API
	muxApi.Handle("/approval-types", loadAzAuthPage(rtApi.CreateApprovalType)).Methods("POST")
	muxApi.Handle("/approval-types/{id}", loadAzAuthPage(rtApi.EditApprovalTypeById)).Methods("PUT")
	muxApi.Handle("/approval-types/{id}/archived", loadAzAuthPage(rtApi.SetIsArchivedApprovalTypeById)).Methods("PUT")
	muxApi.Handle("/approval-types", loadAzAuthPage(rtApi.GetApprovalTypes)).Methods("GET")
	muxApi.Handle("/approval-types/{id}", loadAzAuthPage(rtApi.GetApprovalTypeById)).Methods("GET")

	//EXTERNAL LINKS API
	muxApi.Handle("/external-links", loadAdminPage(rtApi.GetExternalLinks)).Methods("GET")
	muxApi.Handle("/external-links/enabled", loadAzAuthPage(rtApi.GetExternalLinksEnabled)).Methods("GET")
	muxApi.Handle("/external-links/{id}", loadAdminPage(rtApi.GetExternalLinkById)).Methods("GET")
	muxApi.Handle("/external-links", loadAdminPage(rtApi.CreateExternalLinks)).Methods("POST")
	muxApi.Handle("/external-links/{id}", loadAdminPage(rtApi.UpdateExternalLinksById)).Methods("PUT")
	muxApi.Handle("/external-links/{id}", loadAdminPage(rtApi.DeleteExternalLinkById)).Methods("DELETE")

	// APPROVALS API
	muxApi.HandleFunc("/approvals/project/callback", rtApi.UpdateApprovalStatusProjects).Methods("POST")
	muxApi.HandleFunc("/approvals/project/reassign/callback", rtApi.UpdateApprovalReassignApprover)
	muxApi.Handle("/users/{username}/approvals", loadAzAuthPage(rtApi.DownloadProjectApprovalsByUsername))

	// OSS CONTRIBUTION SPONSORS API
	muxApi.Handle("/oss-contribution-sponsors", loadAdminPage((rtApi.GetAllOssContributionSponsors))).Methods("GET")
	muxApi.Handle("/oss-contribution-sponsors/enabled", loadAzAuthPage((rtApi.GetAllEnabledOssContributionSponsors))).Methods("GET")
	muxApi.Handle("/oss-contribution-sponsors", loadAdminPage((rtApi.AddSponsor))).Methods("POST")
	muxApi.Handle("/oss-contribution-sponsors/{id}", loadAdminPage((rtApi.UpdateSponsor))).Methods(("PUT"))

	// LEGACY APIS
	muxApi.Handle("/searchresult/{searchText}", loadGuidAuthApi(rtApi.LegacySearchHandler))
}

func setUtilityRoutes(mux *mux.Router) {
	// UTILITIES
	muxUtility := mux.PathPrefix("/utility").Subrouter()

	// API FOR LOGIC APP
	muxUtility.Handle("/index-org-repos", loadGuidAuthApi(rtApi.IndexOrgRepos)).Methods("GET")
	muxUtility.Handle("/clear-org-repos", loadGuidAuthApi(rtApi.ClearOrgRepos)).Methods("GET")
	muxUtility.Handle("/check-ava-inner-source", loadGuidAuthApi(rtApi.CheckAvaInnerSource)).Methods("GET")
	muxUtility.Handle("/check-ava-open-source", loadGuidAuthApi(rtApi.CheckAvaOpenSource)).Methods("GET")
	muxUtility.Handle("/clear-org-members", loadGuidAuthApi(rtApi.ClearOrgMembers)).Methods("GET")
	muxUtility.Handle("/repo-owner-scan", loadGuidAuthApi(rtApi.RepoOwnerScan)).Methods("GET")
	muxUtility.Handle("/repo-owner-cleanup", loadGuidAuthApi(rtApi.RepoOwnersCleanup)).Methods("GET")
	muxUtility.Handle("/recurring-approval", loadGuidAuthApi(rtApi.RecurringApproval)).Methods("GET")
	muxUtility.Handle("/expiring-invitations", loadGuidAuthApi(rtApi.ExpiringInvitation)).Methods("GET")
	muxUtility.Handle("/fillout-approvers", loadGuidAuthApi(rtApi.FillOutApprovers)).Methods("GET")
	muxUtility.Handle("/fillout-approvalrequest-approvers", loadGuidAuthApi(rtApi.FillOutApprovalRequestApprovers)).Methods("GET")
	muxUtility.Handle("/migrate-oss-sponsors", loadGuidAuthApi(rtApi.MigrateToOssSponsorsTable)).Methods("GET")
}

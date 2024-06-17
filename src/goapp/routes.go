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
	rtOtherRequests "main/routes/pages/otherRequests"
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
	mux.Handle("/repositories", loadAzAuthPage(rtProjects.IndexHandler))
	mux.Handle("/repositories/new", loadAzGHAuthPage(rtProjects.FormHandler))
	mux.Handle("/repositories/view/{githubId}", loadAzAuthPage(rtProjects.ViewByIdHandler))
	mux.Handle("/repositories/view/{org}/{repo}", loadAzAuthPage(rtProjects.ViewHandler))
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

	// OTHER REQUESTS PAGE
	mux.Handle("/other-requests", loadAzGHAuthPage(rtOtherRequests.IndexHandler))
	mux.Handle("/other-requests/organization", loadAzGHAuthPage(rtOtherRequests.RequestNewOrganization))
	mux.Handle("/other-requests/github-copilot", loadAzGHAuthPage(rtOtherRequests.RequestGitHubCopilot))
	mux.Handle("/other-requests/organization-access", loadAzGHAuthPage(rtOtherRequests.RequestOrganizationAccess))

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
	muxApi.Handle("/communities/my", loadAzGHAuthPage(rtApi.GetMyCommunities)).Methods("GET")
	muxApi.Handle("/communities/imanage", loadAzGHAuthPage(rtApi.GetIManageCommunities)).Methods("GET")
	muxApi.Handle("/communities/{id}/members", loadAzAuthPage(rtApi.GetCommunityMembersByCommunityId)).Methods("GET")
	muxApi.Handle("/communities/{id}", loadAzGHAuthPage(rtApi.GetCommunityById)).Methods("GET")
	muxApi.Handle("/communities/isexternal/{isexternal}", loadAzGHAuthPage(rtApi.GetCommunitiesIsExternal)).Methods("GET")
	muxApi.Handle("/communities", loadAzGHAuthPage(rtApi.AddCommunity)).Methods("POST")
	muxApi.Handle("/communities/{id}/members", loadAzGHAuthPage(rtApi.UploadCommunityMembers)).Methods("POST")
	muxApi.Handle("/communities/{id}/sponsors", loadAzGHAuthPage(rtApi.GetCommunitySponsorsByCommunityId)).Methods("GET")
	muxApi.Handle("/communities/{id}/tags", loadAzGHAuthPage(rtApi.GetCommunityTagsByCommunityId)).Methods("GET")
	muxApi.Handle("/communities/{id}/onboarding", loadAzGHAuthPage(rtApi.GetCommunityOnBoardingInfo)).Methods("GET", "POST", "DELETE")
	muxApi.Handle("/communities/{id}/status", loadAzGHAuthPage(rtApi.GetRequestStatusByCommunityId)).Methods("GET")
	muxApi.Handle("/communities/{id}/related-communities", loadAzAuthPage(rtApi.GetRelatedCommunitiesByCommunityId)).Methods("GET")

	// COMMUNITY APPROVERS API
	muxApi.HandleFunc("/community-approvers", rtApi.SubmitCommunityApprover).Methods("POST")
	muxApi.HandleFunc("/community-approvers", rtApi.GetCommunityApproversList).Methods("GET")
	muxApi.HandleFunc("/community-approvers/active", rtApi.GetAllActiveCommunityApprovers).Methods("GET")
	muxApi.HandleFunc("/community-approvers/{id}", rtApi.GetCommunityApproversById).Methods("GET")

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

	// CATEGORY ARTICLES API/
	muxApi.Handle("/categories/{id}/articles", loadAzGHAuthPage(rtApi.GetCategoryArticlesByCategoryId)).Methods("GET")
	muxApi.Handle("/articles/{id}", loadAzGHAuthPage(rtApi.GetCategoryArticlesById)).Methods("GET")
	muxApi.Handle("/articles/{id}", loadAzGHAuthPage(rtApi.UpdateCategoryArticlesById)).Methods("PUT")

	// REPOSITORIES API
	muxApi.Handle("/repositories", loadAzAuthPage(rtApi.GetRepositories)).Methods("GET")
	muxApi.Handle("/repositories/my", loadAzGHAuthPage(rtApi.GetMyRepositories)).Methods("GET")
	muxApi.Handle("/repositories/{id}", loadAzAuthPage(rtApi.GetRepositoriesById)).Methods("GET")
	muxApi.Handle("/repositories/{id}/status", loadAzGHAuthPage(rtApi.GetRequestStatusByRepoId)).Methods("GET")
	muxApi.Handle("/repositories/{orgName}/{repoName}/readme", loadAzAuthPage(rtApi.GetRepositoryReadmeById)).Methods("GET")
	muxApi.Handle("/repositories", loadAzGHAuthPage(rtApi.CreateRepository)).Methods("POST")
	muxApi.Handle("/repositories/{id}/ecattid", loadAzGHAuthPage(rtApi.UpdateRepositoryEcattIdById)).Methods("PUT")

	muxApi.Handle("/repositories/{id}/collaborators", loadAzAuthPage(rtApi.GetRepoCollaboratorsByRepoId)).Methods("GET")
	muxApi.Handle("/repositories/{id}/collaborators/{ghUser}/{permission}", loadAzGHAuthPage(rtApi.AddCollaborator)).Methods("POST")
	muxApi.Handle("/repositories/{id}/collaborators/{ghUser}/{permission}", loadAzGHAuthPage(rtApi.RemoveCollaborator)).Methods("DELETE")

	muxApi.Handle("/repositories/{id}/public", loadAzGHAuthPage(rtApi.RequestMakePublic)).Methods("PUT")
	muxApi.Handle("/repositories/{projectId}/archive/{project}/{organization}/{archive}", loadAzGHAuthPage(rtApi.ArchiveProject)).Methods("PUT")
	muxApi.Handle("/repositories/{projectId}/visibility/{project}/{currentState}/{desiredState}", loadAzGHAuthPage(rtApi.SetVisibility)).Methods("PUT")
	muxApi.Handle("/repositories/{projectId}/transfer", loadAzAuthPage(rtApi.TransferRepository)).Methods("PUT")

	// USERS API
	muxApi.Handle("/users", loadAzAuthPage(rtApi.GetAllUserFromActiveDirectory)).Methods("GET")
	// Retrieve the total number of repositories owned by a me|{user}, categorized by visibility. Default visibility is set to private.
	muxApi.Handle("/users/{user}/repositories/total", loadAzGHAuthPage(rtApi.GetTotalRepositoriesOwnedByUsers)).Methods("GET")
	muxApi.Handle("/users/with-github", loadAzAuthPage(rtApi.GetUsersWithGithub)).Methods("GET")
	muxApi.Handle("/users/{search}/search", loadAzAuthPage(rtApi.SearchUserFromActiveDirectory)).Methods("GET")

	// POPULAR TOPICS API
	muxApi.Handle("/popular-topics", loadAzAuthPage(rtApi.GetPopularTopics)).Methods("GET")

	//APPROVAL TYPES API
	muxApi.Handle("/approval-types", loadAzAuthPage(rtApi.CreateApprovalType)).Methods("POST")
	muxApi.Handle("/approval-types/{id}", loadAzAuthPage(rtApi.EditApprovalTypeById)).Methods("PUT")
	muxApi.Handle("/approval-types/{id}/archived", loadAzAuthPage(rtApi.SetIsArchivedApprovalTypeById)).Methods("PUT")
	muxApi.Handle("/approval-types", loadAzAuthPage(rtApi.GetApprovalTypes)).Methods("GET")
	muxApi.Handle("/approval-types/active", loadAzGHAuthPage(rtApi.GetActiveApprovalTypes)).Methods("GET")
	muxApi.Handle("/approval-types/{id}", loadAzAuthPage(rtApi.GetApprovalTypeById)).Methods("GET")

	//EXTERNAL LINKS API
	muxApi.Handle("/external-links", loadAdminPage(rtApi.GetExternalLinks)).Methods("GET")
	muxApi.Handle("/external-links/enabled", loadAzAuthPage(rtApi.GetExternalLinksEnabled)).Methods("GET")
	muxApi.Handle("/external-links/{id}", loadAdminPage(rtApi.GetExternalLinkById)).Methods("GET")
	muxApi.Handle("/external-links", loadAdminPage(rtApi.CreateExternalLinks)).Methods("POST")
	muxApi.Handle("/external-links/{id}", loadAdminPage(rtApi.UpdateExternalLinksById)).Methods("PUT")
	muxApi.Handle("/external-links/{id}", loadAdminPage(rtApi.DeleteExternalLinkById)).Methods("DELETE")

	// OSS CONTRIBUTION SPONSORS API
	muxApi.Handle("/oss-contribution-sponsors", loadAdminPage((rtApi.GetAllOssContributionSponsors))).Methods("GET")
	muxApi.Handle("/oss-contribution-sponsors/enabled", loadAzAuthPage((rtApi.GetAllEnabledOssContributionSponsors))).Methods("GET")
	muxApi.Handle("/oss-contribution-sponsors", loadAdminPage((rtApi.AddSponsor))).Methods("POST")
	muxApi.Handle("/oss-contribution-sponsors/{id}", loadAdminPage((rtApi.UpdateSponsor))).Methods(("PUT"))

	// OTHER REQUESTS
	muxApi.Handle("/github-organization", loadAzGHAuthPage(rtApi.AddOrganization)).Methods("POST")
	muxApi.Handle("/github-organization", loadAzGHAuthPage(rtApi.GetAllOrganizationRequest)).Methods("GET")
	muxApi.Handle("/github-organization/region", loadAzGHAuthPage(rtApi.GetAllRegionalOrganizations)).Methods("GET")
	muxApi.Handle("/github-organization/{id}/status", loadAzGHAuthPage(rtApi.GetOrganizationApprovalRequests)).Methods("GET")

	muxApi.Handle("/github-copilot", loadAzGHAuthPage(rtApi.AddGitHubCopilot)).Methods("POST")
	muxApi.Handle("/github-copilot", loadAzGHAuthPage(rtApi.GetAllGitHubCopilotRequest)).Methods("GET")
	muxApi.Handle("/github-copilot/{id}/status", loadAzGHAuthPage(rtApi.GetGitHubCopilotApprovalRequests)).Methods("GET")

	muxApi.Handle("/organization-access", loadAzGHAuthPage(rtApi.RequestOrganizationAccess)).Methods("POST")
	muxApi.Handle("/organization-access/me", loadAzGHAuthPage(rtApi.GetMyOrganizationAccess)).Methods("GET")
	muxApi.Handle("/organization-access/{id}/status", loadAzGHAuthPage(rtApi.GetOrganizationAccessApprovalRequests)).Methods("GET")

	//ORGANIZATION APPROVERS API
	muxApi.Handle("/github-organization-approvers/active", loadAzGHAuthPage(rtApi.GetAllActiveOrganizationApprovers)).Methods("GET")

	// APPROVALS API
	muxApi.HandleFunc("/approvals/community/callback", rtApi.UpdateApprovalStatusCommunity).Methods("POST")
	muxApi.HandleFunc("/approvals/organization/callback", rtApi.UpdateApprovalStatusOrganization).Methods("POST")
	muxApi.HandleFunc("/approvals/github-copilot/callback", rtApi.UpdateApprovalStatusCopilot).Methods("POST")
	muxApi.HandleFunc("/approvals/organization-access/callback", rtApi.UpdateApprovalStatusOrganizationAccess).Methods("POST")
	muxApi.HandleFunc("/approvals/community/reassign/callback", rtApi.UpdateCommunityApprovalReassignApprover).Methods("POST")
	muxApi.HandleFunc("/approvals/project/callback", rtApi.UpdateApprovalStatusProjects).Methods("POST")
	muxApi.HandleFunc("/approvals/project/reassign/callback", rtApi.UpdateApprovalReassignApprover).Methods("POST")
	muxApi.Handle("/users/{username}/approvals", loadAzAuthPage(rtApi.DownloadProjectApprovalsByUsername))

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
	muxUtility.Handle("/index-ad-groups", loadGuidAuthApi(rtApi.IndexADGroups)).Methods("GET")
	muxUtility.Handle("/index-regional-organizations", loadGuidAuthApi(rtApi.IndexRegionalOrganizations)).Methods("GET")
}

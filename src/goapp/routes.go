package main

import (
	m "main/middleware"
	ev "main/pkg/envvar"
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
)

func setPageRoutes() {
	httpRouter.NOTFOUND(m.Chain(rtPages.NotFoundHandler, m.AzureAuth()))

	httpRouter.GET("/", m.Chain(rtPages.HomeHandler, m.AzureAuth()))
	httpRouter.GET("/error/ghlogin", m.Chain(rtPages.GHLoginRequire, m.AzureAuth()))

	// SEARCH
	httpRouter.GET("/search", m.Chain(rtSearch.SearchHandler, m.AzureAuth()))

	// ACTIVITIES PAGE
	httpRouter.GET("/activities", m.Chain(rtActivities.IndexHandler, m.AzureAuth()))
	httpRouter.GET("/activities/{action:new}", m.Chain(rtActivities.FormHandler, m.AzureAuth()))
	httpRouter.GET("/activities/{action:edit|view}/{id}", m.Chain(rtActivities.FormHandler, m.AzureAuth()))

	// REPOSITORIES PAGE
	httpRouter.GET("/repositories", m.Chain(rtProjects.IndexHandler, m.AzureAuth()))
	httpRouter.GET("/repositories/new", m.Chain(rtProjects.FormHandler, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/repositories/view/{githubId}", m.Chain(rtProjects.ViewByIdHandler, m.AzureAuth()))
	httpRouter.GET("/repositories/view/{org}/{repo}", m.Chain(rtProjects.ViewHandler, m.AzureAuth()))
	httpRouter.GET("/repositories/makepublic/{id}", m.Chain(rtProjects.MakePublicHandler, m.AzureAuth(), m.GitHubAuth()))

	// GUIDANCE PAGE
	httpRouter.GET("/guidance", m.Chain(rtGuidance.IndexHandler, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/guidance/categories/{id}", m.Chain(rtGuidance.EditCategoryHandler, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/guidance/articles/new", m.Chain(rtGuidance.NewArticleHandler, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/guidance/articles/{id}", m.Chain(rtGuidance.EditArticleHandler, m.AzureAuth(), m.GitHubAuth()))

	// COMMUNITY PAGE
	httpRouter.GET("/communities", m.Chain(rtCommunity.IndexHandler, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/communities/new", m.Chain(rtCommunity.FormHandler, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/communities/{id}", m.Chain(rtCommunity.FormHandler, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/communities/{id}/onboarding", m.Chain(rtCommunity.OnBoardingHandler, m.AzureAuth(), m.GitHubAuth()))

	// OTHER REQUESTS PAGE
	httpRouter.GET("/other-requests", m.Chain(rtOtherRequests.IndexHandler, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/other-requests/organization", m.Chain(rtOtherRequests.RequestNewOrganization, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/other-requests/github-copilot", m.Chain(rtOtherRequests.RequestGitHubCopilot, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/other-requests/organization-access", m.Chain(rtOtherRequests.RequestOrganizationAccess, m.AzureAuth(), m.GitHubAuth()))

	// AUTHENTICATION
	httpRouter.GET("/loginredirect", rtPages.LoginRedirectHandler)
	httpRouter.GET("/gitredirect", rtPages.GitRedirectHandler)
	httpRouter.GET("/clearcookies", rtAzure.ClearCookies)

	// AZURE
	httpRouter.GET("/login/azure", rtAzure.LoginHandler)
	httpRouter.GET("/login/azure/callback", rtAzure.CallbackHandler)
	httpRouter.GET("/logout/azure", rtAzure.LogoutHandler)
	httpRouter.GET("/authentication/azure/inprogress", rtPages.AuthenticationInProgressHandler)
	httpRouter.GET("/authentication/azure/successful", rtPages.AuthenticationSuccessfulHandler)
	httpRouter.GET("/authentication/azure/failed", rtPages.AuthenticationFailedHandler)

	// GITHUB
	httpRouter.GET("/login/github", rtGithub.GithubLoginHandler)
	httpRouter.GET("/login/github/callback", rtGithub.GithubCallbackHandler)
	httpRouter.GET("/login/github/force", rtGithub.GithubForceSaveHandler)
	httpRouter.GET("/logout/github", rtGithub.GitHubLogoutHandler)
	httpRouter.GET("/authentication/github/inprogress", rtPages.GHAuthenticationInProgressHandler)
	httpRouter.GET("/authentication/github/successful", rtPages.AuthenticationSuccessfulHandler)
	httpRouter.GET("/authentication/github/failed", rtPages.AuthenticationFailedHandler)

	// LEGACY REDIRECTS
	httpRouter.GET("/Home/Asset/{assetCode}", rtApi.RedirectAsset)
	httpRouter.GET("/Home/AssetRequestCreation", rtApi.RedirectAssetRequest)
	httpRouter.GET("/Home/AssetRequestCreation/", rtApi.RedirectAssetRequest)
	httpRouter.GET("/Home/Tool/{assetCode}", m.Chain(rtPages.ToolHandler, m.AzureAuth()))
}

func setAdminPageRoutes() {
	// ADMIN
	httpRouter.GET("/admin", m.Chain(rtAdmin.AdminIndexHandler, m.AzureAuth(), m.IsUserAdmin()))

	// COMMUNITY MEMBERS
	httpRouter.GET("/admin/members", m.Chain(rtAdmin.CommunityMembersHandler, m.AzureAuth(), m.IsUserAdmin()))

	// COMMUNITY APPROVERS
	httpRouter.GET("/admin/communityapprovers", m.Chain(rtAdmin.CommunityApproversHandler, m.AzureAuth(), m.IsUserAdmin()))

	// APPROVAL TYPES ADMIN
	httpRouter.GET("/admin/approvaltypes", m.Chain(rtAdmin.ApprovalTypesHandler, m.AzureAuth(), m.IsUserAdmin()))
	httpRouter.GET("/admin/approvaltypes/{action:add}", m.Chain(rtAdmin.ApprovalTypeFormHandler, m.AzureAuth(), m.IsUserAdmin()))
	httpRouter.GET("/admin/approvaltypes/{action:view|edit|delete}/{id}", m.Chain(rtAdmin.ApprovalTypeFormHandler, m.AzureAuth(), m.IsUserAdmin()))

	// CONTRIBUTION AREAS ADMIN
	httpRouter.GET("/admin/contributionareas", m.Chain(rtAdmin.ContributionAreasHandler, m.AzureAuth(), m.IsUserAdmin()))
	httpRouter.GET("/admin/contributionareas/{action:add}", m.Chain(rtAdmin.ContributionAreasFormHandler, m.AzureAuth(), m.IsUserAdmin()))
	httpRouter.GET("/admin/contributionareas/{action:view|edit}/{id}", m.Chain(rtAdmin.ContributionAreasFormHandler, m.AzureAuth(), m.IsUserAdmin()))

	// EXTERNAL LINKS ADMIN
	httpRouter.GET("/admin/externallinks", m.Chain(rtAdmin.ExternalLinksHandler, m.AzureAuth(), m.IsUserAdmin()))
	httpRouter.GET("/admin/externallinks/{action:add}/", m.Chain(rtAdmin.ExternalLinksFormHandler, m.AzureAuth(), m.IsUserAdmin()))
	httpRouter.GET("/admin/externallinks/{action:view|edit}/{id}", m.Chain(rtAdmin.ExternalLinksFormHandler, m.AzureAuth(), m.IsUserAdmin()))

	// OSS CONTRIBUTION SPONSORS ADMIN
	httpRouter.GET("/admin/osscontributionsponsors", m.Chain(rtAdmin.OssContributionSponsorsHandler, m.AzureAuth(), m.IsUserAdmin()))
	httpRouter.GET("/admin/osscontributionsponsors/form", m.Chain(rtAdmin.OssContributionSponsorsFormHandler, m.AzureAuth(), m.IsUserAdmin()))

	// MANAGE ORGANIZATIONS
	httpRouter.GET("/admin/manage-organizations", m.Chain(rtAdmin.RegionalOrganizationHandler, m.AzureAuth(), m.IsUserAdmin()))
}

func setApiRoutes() {
	// APIS

	// ACTIVITY TYPES API
	httpRouter.GET("/api/activity-types", m.Chain(cont.ActivityType.GetActivityTypes, m.AzureAuth(), m.GitHubAuth()))

	// ACTIVITY API
	httpRouter.POST("/api/activities", m.Chain(cont.Activity.CreateActivity, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/activities", m.Chain(cont.Activity.GetActivities, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/activities/{id}", m.Chain(cont.Activity.GetActivityById, m.AzureAuth(), m.GitHubAuth()))

	// COMMUNITIES API
	httpRouter.GET("/api/communities", m.Chain(rtApi.GetCommunities, m.AzureAuth()))
	httpRouter.GET("/api/communities/my", m.Chain(rtApi.GetMyCommunities, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/communities/imanage", m.Chain(rtApi.GetIManageCommunities, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/communities/{id}/members", m.Chain(rtApi.GetCommunityMembersByCommunityId, m.AzureAuth()))
	httpRouter.GET("/api/communities/{id}", m.Chain(rtApi.GetCommunityById, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/communities/isexternal/{isexternal}", m.Chain(rtApi.GetCommunitiesIsExternal, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.POST("/api/communities", m.Chain(rtApi.AddCommunity, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.POST("/api/communities/{id}/members", m.Chain(rtApi.UploadCommunityMembers, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/communities/{id}/sponsors", m.Chain(rtApi.GetCommunitySponsorsByCommunityId, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/communities/{id}/tags", m.Chain(rtApi.GetCommunityTagsByCommunityId, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/communities/{id}/onboarding", m.Chain(rtApi.GetCommunityOnBoardingInfo, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.POST("/api/communities/{id}/onboarding", m.Chain(rtApi.GetCommunityOnBoardingInfo, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.DELETE("/api/communities/{id}/onboarding", m.Chain(rtApi.GetCommunityOnBoardingInfo, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/communities/{id}/status", m.Chain(rtApi.GetRequestStatusByCommunityId, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/communities/{id}/related-communities", m.Chain(rtApi.GetRelatedCommunitiesByCommunityId, m.AzureAuth()))

	// COMMUNITY APPROVERS API
	httpRouter.POST("/api/community-approvers", m.Chain(rtApi.SubmitCommunityApprover, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/community-approvers", m.Chain(rtApi.GetCommunityApproversList, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/community-approvers/active", m.Chain(rtApi.GetAllActiveCommunityApprovers, m.AzureAuth(), m.GitHubAuth()))

	// CONTRIBUTION AREAS API
	httpRouter.POST("/api/contribution-areas", m.Chain(cont.ContributionArea.CreateContributionAreas, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/contribution-areas", m.Chain(cont.ContributionArea.GetContributionAreas, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.PUT("/api/contribution-areas/{id}", m.Chain(cont.ContributionArea.UpdateContributionArea, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/contribution-areas/{id}", m.Chain(cont.ContributionArea.GetContributionAreaById, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/activities/{id}/contribution-areas", m.Chain(rtApi.GetContributionAreasByActivityId, m.AzureAuth(), m.GitHubAuth()))

	// CATEGORIES API
	httpRouter.GET("/api/categories", m.Chain(cont.Category.GetCategories, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/categories/{id}", m.Chain(cont.Category.GetCategoryById, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.PUT("/api/categories/{id}", m.Chain(cont.Category.UpdateCategory, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.POST("/api/categories", m.Chain(cont.Category.CreateCategory, m.AzureAuth(), m.GitHubAuth()))

	// CATEGORY ARTICLES API/
	httpRouter.GET("/api/categories/{id}/articles", m.Chain(cont.Article.GetArticlesByCategoryId, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/articles/{id}", m.Chain(cont.Article.GetArticleById, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.PUT("/api/articles/{id}", m.Chain(cont.Article.UpdateArticle, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.POST("/api/articles/{id}/articles", m.Chain(cont.Article.CreateArticle, m.AzureAuth(), m.GitHubAuth()))
	// REPOSITORIES API
	httpRouter.GET("/api/repositories", m.Chain(rtApi.GetRepositories, m.AzureAuth()))
	httpRouter.GET("/api/repositories/my", m.Chain(rtApi.GetMyRepositories, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/repositories/{id}", m.Chain(rtApi.GetRepositoriesById, m.AzureAuth()))
	httpRouter.GET("/api/repositories/{id}/status", m.Chain(rtApi.GetRequestStatusByRepoId, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/repositories/{orgName}/{repoName}/readme", m.Chain(rtApi.GetRepositoryReadmeById, m.AzureAuth()))
	httpRouter.POST("/api/repositories", m.Chain(rtApi.CreateRepository, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.PUT("/api/repositories/{id}/ecattid", m.Chain(rtApi.UpdateRepositoryEcattIdById, m.AzureAuth(), m.GitHubAuth()))

	httpRouter.GET("/api/repositories/{id}/collaborators", m.Chain(rtApi.GetRepoCollaboratorsByRepoId, m.AzureAuth()))
	httpRouter.POST("/api/repositories/{id}/collaborators/{ghUser}/{permission}", m.Chain(rtApi.AddCollaborator, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.DELETE("/api/repositories/{id}/collaborators/{ghUser}/{permission}", m.Chain(rtApi.RemoveCollaborator, m.AzureAuth(), m.GitHubAuth()))

	httpRouter.PUT("/api/repositories/{id}/public", m.Chain(rtApi.RequestMakePublic, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.PUT("/api/repositories/{projectId}/archive/{project}/{organization}/{archive}", m.Chain(rtApi.ArchiveProject, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.PUT("/api/repositories/{projectId}/visibility/{project}/{currentState}/{desiredState}", m.Chain(rtApi.SetVisibility, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.PUT("/api/repositories/{projectId}/transfer", m.Chain(rtApi.TransferRepository, m.AzureAuth()))

	// USERS API
	httpRouter.GET("/api/users", m.Chain(rtApi.GetAllUserFromActiveDirectory, m.AzureAuth()))
	// Retrieve the total number of repositories owned by a me|{user}, categorized by visibility. Default visibility is set to private.
	httpRouter.GET("/api/users/{user}/repositories/total", m.Chain(rtApi.GetTotalRepositoriesOwnedByUsers, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/users/with-github", m.Chain(rtApi.GetUsersWithGithub, m.AzureAuth()))
	httpRouter.GET("/api/users/{search}/search", m.Chain(rtApi.SearchUserFromActiveDirectory, m.AzureAuth()))

	// POPULAR TOPICS API
	httpRouter.GET("/api/popular-topics", m.Chain(cont.Topic.GetPopularTopics, m.AzureAuth()))

	//APPROVAL TYPES API
	httpRouter.POST("/api/approval-types", m.Chain(cont.ApprovalType.CreateApprovalType, m.AzureAuth()))
	httpRouter.PUT("/api/approval-types/{id}", m.Chain(rtApi.EditApprovalTypeById, m.AzureAuth()))
	httpRouter.PUT("/api/approval-types/{id}/archived", m.Chain(rtApi.SetIsArchivedApprovalTypeById, m.AzureAuth()))
	httpRouter.GET("/api/approval-types", m.Chain(cont.ApprovalType.GetApprovalTypes, m.AzureAuth()))
	httpRouter.GET("/api/approval-types/active", m.Chain(rtApi.GetActiveApprovalTypes, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/approval-types/{id}", m.Chain(cont.ApprovalType.GetApprovalTypeById, m.AzureAuth()))

	// REPOSITORY APPROVERS
	httpRouter.GET("/api/repository-approvers/legal", m.Chain(cont.RepositoryApprover.GetLegalApprovers, m.GuidAuth()))

	//EXTERNAL LINKS API
	httpRouter.GET("/api/external-links", m.Chain(cont.ExternalLink.GetExternalLinks, m.AzureAuth(), m.IsUserAdmin()))
	httpRouter.GET("/api/external-links/enabled", m.Chain(cont.ExternalLink.GetEnabledExternalLinks, m.AzureAuth()))
	httpRouter.GET("/api/external-links/{id}", m.Chain(cont.ExternalLink.GetExternalLinkById, m.AzureAuth(), m.IsUserAdmin()))
	httpRouter.POST("/api/external-links", m.Chain(cont.ExternalLink.CreateExternalLink, m.AzureAuth(), m.IsUserAdmin()))
	httpRouter.PUT("/api/external-links/{id}", m.Chain(cont.ExternalLink.UpdateExternalLinkById, m.AzureAuth(), m.IsUserAdmin()))
	httpRouter.DELETE("/api/external-links/{id}", m.Chain(cont.ExternalLink.RemoveExternalLinkById, m.AzureAuth(), m.IsUserAdmin()))

	// OSS CONTRIBUTION SPONSORS API
	httpRouter.GET("/api/oss-contribution-sponsors", m.Chain(cont.OssContributionSponsor.GetOssContributionSponsors, m.AzureAuth(), m.IsUserAdmin()))
	httpRouter.GET("/api/oss-contribution-sponsors/enabled", m.Chain(cont.OssContributionSponsor.GetEnabledOssContributionSponsors, m.AzureAuth()))
	httpRouter.POST("/api/oss-contribution-sponsors", m.Chain(cont.OssContributionSponsor.CreateOssContributionSponsor, m.AzureAuth(), m.IsUserAdmin()))
	httpRouter.PUT("/api/oss-contribution-sponsors/{id}", m.Chain(cont.OssContributionSponsor.UpdateOssContributionSponsor, m.AzureAuth(), m.IsUserAdmin()))

	// REGIONAL ORGANIZATIONS API
	httpRouter.GET("/api/enterprise-organizations", m.Chain(rtApi.GetEnterpriseOrganizations, m.AzureAuth()))
	httpRouter.GET("/api/regional-organizations", m.Chain(rtApi.GetRegionalOrganizationByOption, m.AzureAuth()))
	httpRouter.GET("/api/regional-organizations/{id}", m.Chain(rtApi.GetRegionalOrganizationById, m.AzureAuth()))
	httpRouter.POST("/api/regional-organizations", m.Chain(rtApi.InsertRegionalOrganization, m.AzureAuth()))
	httpRouter.PUT("/api/regional-organizations/{id}", m.Chain(rtApi.UpdateRegionalOrganization, m.AzureAuth()))

	// OTHER REQUESTS
	httpRouter.POST("/api/github-organization", m.Chain(rtApi.AddOrganization, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/github-organization", m.Chain(rtApi.GetAllOrganizationRequest, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/github-organization/region", m.Chain(rtApi.GetAllRegionalOrganizations, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/github-organization/region/name", m.Chain(rtApi.GetAllRegionalOrganizationsName))
	httpRouter.GET("/api/github-organization/{id}/status", m.Chain(rtApi.GetOrganizationApprovalRequests, m.AzureAuth(), m.GitHubAuth()))

	httpRouter.POST("/api/github-copilot", m.Chain(rtApi.AddGitHubCopilot, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/github-copilot", m.Chain(rtApi.GetAllGitHubCopilotRequest, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/github-copilot/{id}/status", m.Chain(rtApi.GetGitHubCopilotApprovalRequests, m.AzureAuth(), m.GitHubAuth()))

	httpRouter.POST("/api/organization-access", m.Chain(rtApi.RequestOrganizationAccess, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/organization-access/me", m.Chain(rtApi.GetMyOrganizationAccess, m.AzureAuth(), m.GitHubAuth()))
	httpRouter.GET("/api/organization-access/{id}/status", m.Chain(rtApi.GetOrganizationAccessApprovalRequests, m.AzureAuth(), m.GitHubAuth()))

	//ORGANIZATION APPROVERS API
	httpRouter.GET("/api/github-organization-approvers/active", m.Chain(rtApi.GetAllActiveOrganizationApprovers, m.AzureAuth(), m.GitHubAuth()))

	// APPROVALS API
	httpRouter.POST("/api/approvals/community/callback", m.Chain(rtApi.UpdateApprovalStatusCommunity, m.GuidAuth()))
	httpRouter.POST("/api/approvals/organization/callback", m.Chain(rtApi.UpdateApprovalStatusOrganization, m.GuidAuth()))
	httpRouter.POST("/api/approvals/github-copilot/callback", m.Chain(rtApi.UpdateApprovalStatusCopilot, m.GuidAuth()))
	httpRouter.POST("/api/approvals/organization-access/callback", m.Chain(rtApi.UpdateApprovalStatusOrganizationAccess, m.GuidAuth()))
	httpRouter.POST("/api/approvals/community/reassign/callback", m.Chain(rtApi.UpdateCommunityApprovalReassignApprover, m.GuidAuth()))
	httpRouter.POST("/api/approvals/project/callback", m.Chain(rtApi.UpdateApprovalStatusProjects, m.GuidAuth()))
	httpRouter.POST("/api/approvals/project/reassign/callback", m.Chain(rtApi.UpdateApprovalReassignApprover, m.GuidAuth()))
	httpRouter.GET("/api/users/{username}/approvals", m.Chain(rtApi.DownloadProjectApprovalsByUsername, m.GuidAuth()))

	// LEGACY APIS
	httpRouter.GET("/api/searchresult/{searchText}", m.Chain(rtApi.LegacySearchHandler, m.GuidAuth()))

	// SEARCH API
	httpRouter.GET("/api/search", m.Chain(rtSearch.GetSearchResults, m.AzureAuth()))
}

func setUtilityRoutes() {
	// UTILITIES

	// API FOR LOGIC APP
	httpRouter.GET("/utility/index-org-repos", m.Chain(rtApi.IndexOrgRepos, m.GuidAuth()))
	httpRouter.GET("/utility/clear-org-repos", m.Chain(rtApi.ClearOrgRepos, m.GuidAuth()))
	httpRouter.GET("/utility/check-ava-inner-source", m.Chain(rtApi.CheckAvaInnerSource, m.GuidAuth()))
	httpRouter.GET("/utility/check-ava-open-source", m.Chain(rtApi.CheckAvaOpenSource, m.GuidAuth()))
	httpRouter.GET("/utility/clear-org-members", m.Chain(rtApi.ClearOrgMembers, m.GuidAuth()))
	httpRouter.GET("/utility/demote-outside-collaborators-admin", m.Chain(rtApi.DemoteOutsideCollaboratorAdmin, m.GuidAuth()))
	httpRouter.GET("/utility/repo-owner-scan", m.Chain(rtApi.RepoOwnerScan, m.GuidAuth()))
	httpRouter.GET("/utility/repo-owner-cleanup", m.Chain(rtApi.RepoOwnersCleanup, m.GuidAuth()))
	httpRouter.GET("/utility/recurring-approval", m.Chain(rtApi.RecurringApproval, m.GuidAuth()))
	httpRouter.GET("/utility/expiring-invitations", m.Chain(rtApi.ExpiringInvitation, m.GuidAuth()))
	httpRouter.GET("/utility/index-ad-groups", m.Chain(rtApi.IndexADGroups, m.GuidAuth()))
	httpRouter.GET("/utility/index-regional-organizations", m.Chain(rtApi.IndexRegionalOrganizations, m.GuidAuth()))
	httpRouter.GET("/utility/scan-community-organizations", m.Chain(rtApi.ScanCommunityOrganizations, m.GuidAuth()))
}

func serve() {
	port := ev.GetEnvVar("PORT", "8080")
	httpRouter.SERVE(port)
}

-- This file contains SQL statements that will be executed after the build script.
/* INITIAL DATA FOR APPROVAL STATUS */
SET IDENTITY_INSERT ApprovalStatus ON IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 1
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (1, 'New') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 2
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (2, 'InReview') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 3
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (3, 'Rejected') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 4
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (4, 'NonCompliant') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 5
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (5, 'Approved') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 6
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (6, 'Retired') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalStatus]
        WHERE [Id] = 7
    )
INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name])
VALUES (7, 'Archived')
SET IDENTITY_INSERT ApprovalStatus OFF
    /* INITIAL DATA FOR APPROVAL TYPES */
SET IDENTITY_INSERT ApprovalTypes ON IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalTypes]
        WHERE [Id] = 1
    )
INSERT INTO [dbo].[ApprovalTypes] ([Id], [Name])
VALUES (1, 'Intellectual Property') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalTypes]
        WHERE [Id] = 2
    )
INSERT INTO [dbo].[ApprovalTypes] ([Id], [Name])
VALUES (2, 'Legal') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[ApprovalTypes]
        WHERE [Id] = 3
    )
INSERT INTO [dbo].[ApprovalTypes] ([Id], [Name])
VALUES (3, 'Security')
SET IDENTITY_INSERT ApprovalTypes OFF
SET IDENTITY_INSERT Visibility ON IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[Visibility]
        WHERE [Id] = 1
    )
INSERT INTO [dbo].[Visibility] ([Id], [Name])
VALUES (1, 'Private') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[Visibility]
        WHERE [Id] = 2
    )
INSERT INTO [dbo].[Visibility] ([Id], [Name])
VALUES (2, 'Internal') IF NOT EXISTS (
        SELECT [Id]
        FROM [dbo].[Visibility]
        WHERE [Id] = 3
    )
INSERT INTO [dbo].[Visibility] ([Id], [Name])
VALUES (3, 'Public')
SET IDENTITY_INSERT Visibility OFF

/* RENAME ALL TABLES */

/* ActivityTypes > ActivityType */
EXEC sp_rename 'dbo.ActivityTypes', 'ActivityType'

/* ApprovalRequestApprovers > ApprovalRequestApprover */
EXEC sp_rename 'dbo.ApprovalRequestApprovers', 'ApprovalRequestApprover'

/* ApprovalStatus > ApprovalStatus */


/* ApprovalTypes > RepositoryApprovalType */
EXEC sp_rename 'dbo.ApprovalTypes', 'RepositoryApprovalType'

/* Approvers > RepositoryApprover */
EXEC sp_rename 'dbo.Approvers', 'RepositoryApprover'

/* Category > GuidanceCategory */
EXEC sp_rename 'dbo.Category', 'GuidanceCategory'

/* CategoryArticles > GuidanceCategoryArticle */
EXEC sp_rename 'dbo.CategoryArticles', 'GuidanceCategoryArticle'

/* Communities > Community */
EXEC sp_rename 'dbo.Communities', 'Community'

/* CommunityActivities > CommunityActivity */
EXEC sp_rename 'dbo.CommunityActivities', 'CommunityActivity'

/* CommunityActivitiesContributionAreas > CommunityActivityContributionArea */
EXEC sp_rename 'dbo.CommunityActivitiesContributionAreas', 'CommunityActivityContributionArea'

/* CommunityActivitiesHelpTypes > CommunityActivityHelpType */
EXEC sp_rename 'dbo.CommunityActivitiesHelpTypes', 'CommunityActivityHelpType'

/* CommunityApprovalRequests > CommunityApprovalRequest */
EXEC sp_rename 'dbo.CommunityApprovalRequests', 'CommunityApprovalRequest'

/* CommunityApprovals > ApprovalRequest */
EXEC sp_rename 'dbo.CommunityApprovals', 'ApprovalRequest'

/* CommunityApproversList > CommunityApprover */
EXEC sp_rename 'dbo.CommunityApproversList', 'CommunityApprover'

/* CommunityMembers > CommunityMember */
EXEC sp_rename 'dbo.CommunityMembers', 'CommunityMember'

/* CommunitySponsors > CommunitySponsor */
EXEC sp_rename 'dbo.CommunitySponsors', 'CommunitySponsor'

/* CommunityTags > CommunityTag */
EXEC sp_rename 'dbo.CommunityTags', 'CommunityTag'

/* ContributionAreas > ContributionArea */
EXEC sp_rename 'dbo.ContributionAreas', 'ContributionArea'

/* ExternalLinks > ExternalLink */
EXEC sp_rename 'dbo.ExternalLinks', 'ExternalLink'

/* GitHubAccess > GitHubAccessDirectoryGroup */
EXEC sp_rename 'dbo.GitHubAccess', 'GitHubAccessDirectoryGroup'

/* GitHubCopilot > GitHubCopilot */


/* GitHubCopilotApprovalRequests > GitHubCopilotApprovalRequest */
EXEC sp_rename 'dbo.GitHubCopilotApprovalRequests', 'GitHubCopilotApprovalRequest'

/* OrganizationAccess > OrganizationAccess */


/* OrganizationAccessApprovalRequests > OrganizationAccessApprovalRequest */
EXEC sp_rename 'dbo.OrganizationAccessApprovalRequests', 'OrganizationAccessApprovalRequest'

/* OrganizationApprovalRequests > OrganizationApprovalRequest */
EXEC sp_rename 'dbo.OrganizationApprovalRequests', 'OrganizationApprovalRequest'

/* Organizations > Organization */
EXEC sp_rename 'dbo.Organizations', 'Organization'

/* OSSContributionSponsors > OSSContributionSponsor */
EXEC sp_rename 'dbo.OSSContributionSponsors', 'OSSContributionSponsor'

/* ProjectApprovals > RepositoryApproval */
EXEC sp_rename 'dbo.ProjectApprovals', 'RepositoryApproval'

/* Projects > Repository */
EXEC sp_rename 'dbo.Projects', 'Repository'

/* RepoTopics > RepositoryTopic */
EXEC sp_rename 'dbo.RepoTopics', 'RepositoryTopic'

/* RegionalOrganizations > RegionalOrganization */
EXEC sp_rename 'dbo.RegionalOrganizations', 'RegionalOrganization'

/* RelatedCommunities > RelatedCommunity */
EXEC sp_rename 'dbo.RelatedCommunities', 'RelatedCommunity'

/* RepoOwners > RepositoryOwner */
EXEC sp_rename 'dbo.RepoOwners', 'RepositoryOwner'

/* Users > User */
EXEC sp_rename 'dbo.Users', 'User'

/* Visibility > Visibility */


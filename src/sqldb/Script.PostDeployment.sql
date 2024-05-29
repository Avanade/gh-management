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
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ActivityType'))
BEGIN
    EXEC sp_rename 'dbo.ActivityTypes', 'ActivityType'
END

/* ApprovalRequestApprovers > ApprovalRequestApprover */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ApprovalRequestApprover'))
BEGIN
    EXEC sp_rename 'dbo.ApprovalRequestApprovers', 'ApprovalRequestApprover'
END

/* ApprovalStatus > ApprovalStatus */


/* ApprovalTypes > RepositoryApprovalType */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'RepositoryApprovalType'))
BEGIN
    EXEC sp_rename 'dbo.ApprovalTypes', 'RepositoryApprovalType'
END

/* Approvers > RepositoryApprover */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'RepositoryApprover'))
BEGIN
    EXEC sp_rename 'dbo.Approvers', 'RepositoryApprover'
END

/* Category > GuidanceCategory */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'GuidanceCategory'))
BEGIN
    EXEC sp_rename 'dbo.Category', 'GuidanceCategory'
END

/* CategoryArticles > GuidanceCategoryArticle */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'GuidanceCategoryArticle'))
BEGIN
    EXEC sp_rename 'dbo.CategoryArticles', 'GuidanceCategoryArticle'
END

/* Communities > Community */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'Community'))
BEGIN
    EXEC sp_rename 'dbo.Communities', 'Community'
END

/* CommunityActivities > CommunityActivity */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityActivity'))
BEGIN
    EXEC sp_rename 'dbo.CommunityActivities', 'CommunityActivity'
END

/* CommunityActivitiesContributionAreas > CommunityActivityContributionArea */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityActivityContributionArea'))
BEGIN
    EXEC sp_rename 'dbo.CommunityActivitiesContributionAreas', 'CommunityActivityContributionArea'
END

/* CommunityActivitiesHelpTypes > CommunityActivityHelpType */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityActivityHelpType'))
BEGIN
    EXEC sp_rename 'dbo.CommunityActivitiesHelpTypes', 'CommunityActivityHelpType'
END

/* CommunityApprovalRequests > CommunityApprovalRequest */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityApprovalRequest'))
BEGIN
    EXEC sp_rename 'dbo.CommunityApprovalRequests', 'CommunityApprovalRequest'
END

/* CommunityApprovals > ApprovalRequest */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ApprovalRequest'))
BEGIN
    EXEC sp_rename 'dbo.CommunityApprovals', 'ApprovalRequest'
END

/* CommunityApproversList > CommunityApprover */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityApprover'))
BEGIN
    EXEC sp_rename 'dbo.CommunityApproversList', 'CommunityApprover'
END

/* CommunityMembers > CommunityMember */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityMember'))
BEGIN
    EXEC sp_rename 'dbo.CommunityMembers', 'CommunityMember'
END

/* CommunitySponsors > CommunitySponsor */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunitySponsor'))
BEGIN
    EXEC sp_rename 'dbo.CommunitySponsors', 'CommunitySponsor'
END

/* CommunityTags > CommunityTag */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityTag'))
BEGIN
    EXEC sp_rename 'dbo.CommunityTags', 'CommunityTag'
END

/* ContributionAreas > ContributionArea */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ContributionArea'))
BEGIN
    EXEC sp_rename 'dbo.ContributionAreas', 'ContributionArea'
END

/* ExternalLinks > ExternalLink */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ExternalLink'))
BEGIN
    EXEC sp_rename 'dbo.ExternalLinks', 'ExternalLink'
END

/* GitHubAccess > GitHubAccessDirectoryGroup */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'GitHubAccessDirectoryGroup'))
BEGIN
    EXEC sp_rename 'dbo.GitHubAccess', 'GitHubAccessDirectoryGroup'
END

/* GitHubCopilot > GitHubCopilot */


/* GitHubCopilotApprovalRequests > GitHubCopilotApprovalRequest */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'GitHubCopilotApprovalRequest'))
BEGIN
    EXEC sp_rename 'dbo.GitHubCopilotApprovalRequests', 'GitHubCopilotApprovalRequest'
END

/* OrganizationAccess > OrganizationAccess */


/* OrganizationAccessApprovalRequests > OrganizationAccessApprovalRequest */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'OrganizationAccessApprovalRequest'))
BEGIN
    EXEC sp_rename 'dbo.OrganizationAccessApprovalRequests', 'OrganizationAccessApprovalRequest'
END

/* OrganizationApprovalRequests > OrganizationApprovalRequest */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'OrganizationApprovalRequest'))
BEGIN
    EXEC sp_rename 'dbo.OrganizationApprovalRequests', 'OrganizationApprovalRequest'
END

/* Organizations > Organization */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'Organization'))
BEGIN
    EXEC sp_rename 'dbo.Organizations', 'Organization'
END

/* OSSContributionSponsors > OSSContributionSponsor */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'OSSContributionSponsor'))
BEGIN
    EXEC sp_rename 'dbo.OSSContributionSponsors', 'OSSContributionSponsor'
END

/* ProjectApprovals > RepositoryApproval */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'RepositoryApproval'))
BEGIN
    EXEC sp_rename 'dbo.ProjectApprovals', 'RepositoryApproval'
END

/* Projects > Repository */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'Repository'))
BEGIN
    EXEC sp_rename 'dbo.Projects', 'Repository'
END

/* RepoTopics > RepositoryTopic */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'RepositoryTopic'))
BEGIN
    EXEC sp_rename 'dbo.RepoTopics', 'RepositoryTopic'
END

/* RegionalOrganizations > RegionalOrganization */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'RegionalOrganization'))
BEGIN
    EXEC sp_rename 'dbo.RegionalOrganizations', 'RegionalOrganization'
END

/* RelatedCommunities > RelatedCommunity */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'RelatedCommunity'))
BEGIN
    EXEC sp_rename 'dbo.RelatedCommunities', 'RelatedCommunity'
END

/* RepoOwners > RepositoryOwner */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'RepositoryOwner'))
BEGIN
    EXEC sp_rename 'dbo.RepoOwners', 'RepositoryOwner'
END

/* Users > User */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'User'))
BEGIN
    EXEC sp_rename 'dbo.Users', 'User'
END

/* Visibility > Visibility */


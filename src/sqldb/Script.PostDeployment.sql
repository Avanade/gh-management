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
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ActivityTypes'))
BEGIN
    EXEC sp_rename 'dbo.ActivityTypes', 'ActivityType'
END

/* ApprovalRequestApprovers > ApprovalRequestApprover */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ApprovalRequestApprovers'))
BEGIN
    EXEC sp_rename 'dbo.ApprovalRequestApprovers', 'ApprovalRequestApprover'
END

/* ApprovalStatus > ApprovalStatus */


/* ApprovalTypes > RepositoryApprovalType */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ApprovalTypes'))
BEGIN
    EXEC sp_rename 'dbo.ApprovalTypes', 'RepositoryApprovalType'
END

/* Approvers > RepositoryApprover */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'Approvers'))
BEGIN
    EXEC sp_rename 'dbo.Approvers', 'RepositoryApprover'
END

/* Category > GuidanceCategory */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'Category'))
BEGIN
    EXEC sp_rename 'dbo.Category', 'GuidanceCategory'
END

/* CategoryArticles > GuidanceCategoryArticle */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CategoryArticles'))
BEGIN
    EXEC sp_rename 'dbo.CategoryArticles', 'GuidanceCategoryArticle'
END

/* Communities > Community */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'Communities'))
BEGIN
    EXEC sp_rename 'dbo.Communities', 'Community'
END

/* CommunityActivities > CommunityActivity */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityActivities'))
BEGIN
    EXEC sp_rename 'dbo.CommunityActivities', 'CommunityActivity'
END

/* CommunityActivitiesContributionAreas > CommunityActivityContributionArea */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityActivitiesContributionAreas'))
BEGIN
    EXEC sp_rename 'dbo.CommunityActivitiesContributionAreas', 'CommunityActivityContributionArea'
END

/* CommunityActivitiesHelpTypes > CommunityActivityHelpType */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityActivitiesHelpTypes'))
BEGIN
    EXEC sp_rename 'dbo.CommunityActivitiesHelpTypes', 'CommunityActivityHelpType'
END

/* CommunityApprovalRequests > CommunityApprovalRequest */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityApprovalRequests'))
BEGIN
    EXEC sp_rename 'dbo.CommunityApprovalRequests', 'CommunityApprovalRequest'
END

/* CommunityApprovals > ApprovalRequest */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityApprovals'))
BEGIN
    EXEC sp_rename 'dbo.CommunityApprovals', 'ApprovalRequest'
END

/* CommunityApproversList > CommunityApprover */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityApproversList'))
BEGIN
    EXEC sp_rename 'dbo.CommunityApproversList', 'CommunityApprover'
END

/* CommunityMembers > CommunityMember */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityMembers'))
BEGIN
    EXEC sp_rename 'dbo.CommunityMembers', 'CommunityMember'
END

/* CommunitySponsors > CommunitySponsor */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunitySponsors'))
BEGIN
    EXEC sp_rename 'dbo.CommunitySponsors', 'CommunitySponsor'
END

/* CommunityTags > CommunityTag */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'CommunityTags'))
BEGIN
    EXEC sp_rename 'dbo.CommunityTags', 'CommunityTag'
END

/* ContributionAreas > ContributionArea */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ContributionAreas'))
BEGIN
    EXEC sp_rename 'dbo.ContributionAreas', 'ContributionArea'
END

/* ExternalLinks > ExternalLink */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ExternalLinks'))
BEGIN
    EXEC sp_rename 'dbo.ExternalLinks', 'ExternalLink'
END

/* GitHubAccess > GitHubAccessDirectoryGroup */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'GitHubAccess'))
BEGIN
    EXEC sp_rename 'dbo.GitHubAccess', 'GitHubAccessDirectoryGroup'
END

/* GitHubCopilot > GitHubCopilot */


/* GitHubCopilotApprovalRequests > GitHubCopilotApprovalRequest */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'GitHubCopilotApprovalRequests'))
BEGIN
    EXEC sp_rename 'dbo.GitHubCopilotApprovalRequests', 'GitHubCopilotApprovalRequest'
END

/* OrganizationAccess > OrganizationAccess */


/* OrganizationAccessApprovalRequests > OrganizationAccessApprovalRequest */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'OrganizationAccessApprovalRequests'))
BEGIN
    EXEC sp_rename 'dbo.OrganizationAccessApprovalRequests', 'OrganizationAccessApprovalRequest'
END

/* OrganizationApprovalRequests > OrganizationApprovalRequest */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'OrganizationApprovalRequests'))
BEGIN
    EXEC sp_rename 'dbo.OrganizationApprovalRequests', 'OrganizationApprovalRequest'
END

/* Organizations > Organization */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'Organizations'))
BEGIN
    EXEC sp_rename 'dbo.Organizations', 'Organization'
END

/* OSSContributionSponsors > OSSContributionSponsor */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'OSSContributionSponsors'))
BEGIN
    EXEC sp_rename 'dbo.OSSContributionSponsors', 'OSSContributionSponsor'
END

/* ProjectApprovals > RepositoryApproval */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ProjectApprovals'))
BEGIN
    EXEC sp_rename 'dbo.ProjectApprovals', 'RepositoryApproval'
END

/* Projects > Repository */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'Projects'))
BEGIN
    EXEC sp_rename 'dbo.Projects', 'Repository'
END

/* RepoTopics > RepositoryTopic */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'RepoTopics'))
BEGIN
    EXEC sp_rename 'dbo.RepoTopics', 'RepositoryTopic'
END

/* RegionalOrganizations > RegionalOrganization */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'RegionalOrganizations'))
BEGIN
    EXEC sp_rename 'dbo.RegionalOrganizations', 'RegionalOrganization'
END

/* RelatedCommunities > RelatedCommunity */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'RelatedCommunities'))
BEGIN
    EXEC sp_rename 'dbo.RelatedCommunities', 'RelatedCommunity'
END

/* RepoOwners > RepositoryOwner */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'RepoOwners'))
BEGIN
    EXEC sp_rename 'dbo.RepoOwners', 'RepositoryOwner'
END

/* Users > User */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'Users'))
BEGIN
    EXEC sp_rename 'dbo.Users', 'User'
END

/* Visibility > Visibility */


-- This file contains SQL statements that will be executed after the build script.

/* INITIAL DATA FOR APPROVAL STATUS */
SET IDENTITY_INSERT ApprovalStatus ON 
    IF NOT EXISTS (SELECT [Id] FROM [dbo].[ApprovalStatus] WHERE [Id] = 1)
        INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name]) VALUES (1, 'New') 
    IF NOT EXISTS (SELECT [Id] FROM [dbo].[ApprovalStatus] WHERE [Id] = 2)
        INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name]) VALUES (2, 'InReview') 
    IF NOT EXISTS (SELECT [Id] FROM [dbo].[ApprovalStatus] WHERE [Id] = 3)
        INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name]) VALUES (3, 'Rejected') 
    IF NOT EXISTS (SELECT [Id] FROM [dbo].[ApprovalStatus] WHERE [Id] = 4)
        INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name]) VALUES (4, 'NonCompliant') 
    IF NOT EXISTS (SELECT [Id] FROM [dbo].[ApprovalStatus] WHERE [Id] = 5)
        INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name]) VALUES (5, 'Approved') 
    IF NOT EXISTS (SELECT [Id] FROM [dbo].[ApprovalStatus] WHERE [Id] = 6)
        INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name]) VALUES (6, 'Retired') 
    IF NOT EXISTS (SELECT [Id] FROM [dbo].[ApprovalStatus] WHERE [Id] = 7)
        INSERT INTO [dbo].[ApprovalStatus] ([Id], [Name]) VALUES (7, 'Archived')
SET IDENTITY_INSERT ApprovalStatus OFF

/* INITIAL DATA FOR APPROVAL TYPES */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'RepositoryApprovalType'))
    BEGIN
        SET IDENTITY_INSERT RepositoryApprovalType ON
            IF NOT EXISTS (SELECT [Id] FROM [dbo].[RepositoryApprovalType] WHERE [Id] = 1)
                INSERT INTO [dbo].[RepositoryApprovalType] ([Id], [Name]) VALUES (1, 'Intellectual Property') 
            IF NOT EXISTS (SELECT [Id] FROM [dbo].[RepositoryApprovalType] WHERE [Id] = 2)
                INSERT INTO [dbo].[RepositoryApprovalType] ([Id], [Name]) VALUES (2, 'Legal') 
            IF NOT EXISTS (SELECT [Id] FROM [dbo].[RepositoryApprovalType] WHERE [Id] = 3)
                INSERT INTO [dbo].[RepositoryApprovalType] ([Id], [Name]) VALUES (3, 'Security')
        SET IDENTITY_INSERT RepositoryApprovalType OFF
    END
ELSE
    BEGIN
        SET IDENTITY_INSERT ApprovalTypes ON
            IF NOT EXISTS (SELECT [Id] FROM [dbo].[ApprovalTypes] WHERE [Id] = 1)
                INSERT INTO [dbo].[ApprovalTypes] ([Id], [Name]) VALUES (1, 'Intellectual Property') 
            IF NOT EXISTS (SELECT [Id] FROM [dbo].[ApprovalTypes] WHERE [Id] = 2)
                INSERT INTO [dbo].[ApprovalTypes] ([Id], [Name]) VALUES (2, 'Legal') 
            IF NOT EXISTS (SELECT [Id] FROM [dbo].[ApprovalTypes] WHERE [Id] = 3)
                INSERT INTO [dbo].[ApprovalTypes] ([Id], [Name]) VALUES (3, 'Security')
        SET IDENTITY_INSERT ApprovalTypes OFF
    END

/* INITIAL DATA FOR VISIBILITY */
SET IDENTITY_INSERT Visibility ON 
    IF NOT EXISTS (SELECT [Id] FROM [dbo].[Visibility] WHERE [Id] = 1)
        INSERT INTO [dbo].[Visibility] ([Id], [Name]) VALUES (1, 'Private') 
    IF NOT EXISTS (SELECT [Id] FROM [dbo].[Visibility] WHERE [Id] = 2)
        INSERT INTO [dbo].[Visibility] ([Id], [Name]) VALUES (2, 'Internal') 
    IF NOT EXISTS (SELECT [Id] FROM [dbo].[Visibility] WHERE [Id] = 3)
        INSERT INTO [dbo].[Visibility] ([Id], [Name]) VALUES (3, 'Public')
SET IDENTITY_INSERT Visibility OFF

/*  DROP ALL CONSTRAINT */
IF (NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ActivityType'))
BEGIN
    -- ApprovalRequestApprovers
    ALTER TABLE [dbo].[ApprovalRequestApprovers] DROP CONSTRAINT [PK_ApprovalRequestApprover]
    ALTER TABLE [dbo].[ApprovalRequestApprovers] DROP CONSTRAINT [FK_ApprovalRequestApprover_ProjectApprovals]
    ALTER TABLE [dbo].[ApprovalRequestApprovers] DROP CONSTRAINT [FK_ApprovalRequestApprover_Users]

    -- ApprovalTypes
    ALTER TABLE [dbo].[ApprovalTypes] DROP CONSTRAINT [FK_ApprovalTypes_Users]

    -- Approvers
    ALTER TABLE [dbo].[Approvers] DROP CONSTRAINT [PK_Approver]
    ALTER TABLE [dbo].[Approvers] DROP CONSTRAINT [FK_Approvers_ApprovalTypes]
    ALTER TABLE [dbo].[Approvers] DROP CONSTRAINT [FK_Approvers_Users]

    -- CategoryArticles
    ALTER TABLE [dbo].[CategoryArticles] DROP CONSTRAINT [FK_CategoryArticles_Category]

    -- Communities
    ALTER TABLE [dbo].[Communities] DROP CONSTRAINT [FK_Communities_ApprovalStatus] 

    -- CommunityActivities
    ALTER TABLE [dbo].[CommunityActivities] DROP CONSTRAINT [FK_CommunityActivities_Communities]
    ALTER TABLE [dbo].[CommunityActivities] DROP CONSTRAINT [FK_CommunityActivities_ActivityTypes]

    -- CommunityActivitiesContributionAreas
    ALTER TABLE [dbo].[CommunityActivitiesContributionAreas] DROP CONSTRAINT [FK_CommunityActivitiesCA_CommunityActivity]
    ALTER TABLE [dbo].[CommunityActivitiesContributionAreas] DROP CONSTRAINT [FK_CommunityActivitiesCA_ContributionAreas]

    -- CommunityApprovalRequests
    ALTER TABLE [dbo].[CommunityApprovalRequests] DROP CONSTRAINT [PK_CommunityApprovalRequests]
    ALTER TABLE [dbo].[CommunityApprovalRequests] DROP CONSTRAINT [FK_CommunityApprovalRequests_Communities]
    ALTER TABLE [dbo].[CommunityApprovalRequests] DROP CONSTRAINT [FK_CommunityApprovalRequests_CommunityApprovals]

    -- CommunityApprovals
    ALTER TABLE [dbo].[CommunityApprovals] DROP CONSTRAINT [FK_CommunityApprovals_Users]
    ALTER TABLE [dbo].[CommunityApprovals] DROP CONSTRAINT [FK_CommunityApprovals_ApprovalStatus]

    -- CommunityApproversList
    ALTER TABLE [dbo].[CommunityApproversList] DROP CONSTRAINT [AK_ApproverUserPrincipalName_Categoy]
    ALTER TABLE [dbo].[CommunityApproversList] DROP CONSTRAINT [FK_CommunityApproversList_Users]

    -- CommunityMembers
    ALTER TABLE [dbo].[CommunityMembers] DROP CONSTRAINT [PK_CommunityMembers]
    ALTER TABLE [dbo].[CommunityMembers] DROP CONSTRAINT [FK_CommunityMembers_Communities]

    -- CommunitySponsors
    ALTER TABLE [dbo].[CommunitySponsors] DROP CONSTRAINT [AK_CommunityId_UserPrincipalName]
    ALTER TABLE [dbo].[CommunitySponsors] DROP CONSTRAINT [FK_CommunitySponsors_Communities]
    ALTER TABLE [dbo].[CommunitySponsors] DROP CONSTRAINT [FK_CommunitySponsors_Users]

    -- CommunityTags
    ALTER TABLE [dbo].[CommunityTags] DROP CONSTRAINT [AK_CommunityId_Tag]
    ALTER TABLE [dbo].[CommunityTags] DROP CONSTRAINT [FK_CommunityTags_Communities]

    -- GitHubCopilot
    ALTER TABLE [dbo].[GitHubCopilot] DROP CONSTRAINT [FK_GitHubCopilot_RegionalOrganizations]

    -- GitHubCopilotApprovalRequests
    ALTER TABLE [dbo].[GitHubCopilotApprovalRequests] DROP CONSTRAINT [PK_GitHubCopilotApprovalRequests]
    ALTER TABLE [dbo].[GitHubCopilotApprovalRequests] DROP CONSTRAINT [FK_GitHubCopilotApprovalRequests_GitHubCopilot]
    ALTER TABLE [dbo].[GitHubCopilotApprovalRequests] DROP CONSTRAINT [FK_GitHubCopilotApprovalRequests_CommunityApprovals]

    -- OrganizationAccess
    ALTER TABLE [dbo].[OrganizationAccess] DROP CONSTRAINT [FK_OrganizationAccess_Users]
    ALTER TABLE [dbo].[OrganizationAccess] DROP CONSTRAINT [FK_OrganizationAccess_RegionalOrganizations]

    -- OrganizationAccessApprovalRequests
    ALTER TABLE [dbo].[OrganizationAccessApprovalRequests] DROP CONSTRAINT [PK_OrganizationAccessApprovalRequests]
    ALTER TABLE [dbo].[OrganizationAccessApprovalRequests] DROP CONSTRAINT [FK_OrganizationAccessApprovalRequests_OrganizationAccess]
    ALTER TABLE [dbo].[OrganizationAccessApprovalRequests] DROP CONSTRAINT [FK_OrganizationAccessApprovalRequests_CommunityApprovals]

    -- OrganizationApprovalRequests
    ALTER TABLE [dbo].[OrganizationApprovalRequests] DROP CONSTRAINT [PK_OrganizationApprovalRequests]
    ALTER TABLE [dbo].[OrganizationApprovalRequests] DROP CONSTRAINT [FK_OrganizationApprovalRequests_Organizations]
    ALTER TABLE [dbo].[OrganizationApprovalRequests] DROP CONSTRAINT [FK_OrganizationApprovalRequests_CommunityApprovals]

    -- Organizations
    ALTER TABLE [dbo].[Organizations] DROP CONSTRAINT [FK_Organizations_RegionalOrganizations]

    -- ProjectApprovals
    ALTER TABLE [dbo].[ProjectApprovals] DROP CONSTRAINT [FK_ProjectApprovals_Projects]
    ALTER TABLE [dbo].[ProjectApprovals] DROP CONSTRAINT [FK_ProjectApprovals_ApprovalTypes]
    ALTER TABLE [dbo].[ProjectApprovals] DROP CONSTRAINT [FK_ProjectApprovals_ApprovalStatus]

    -- Projects
    ALTER TABLE [dbo].[Projects] DROP CONSTRAINT [FK_Projects_ApprovalStatus]
    ALTER TABLE [dbo].[Projects] DROP CONSTRAINT [FK_Projects_Visibility]
    ALTER TABLE [dbo].[Projects] DROP CONSTRAINT [FK_Projects_OSSContributionSponsors]

    -- RepoTopics
    ALTER TABLE [dbo].[RepoTopics] DROP CONSTRAINT [PK_RepoTopics]
    ALTER TABLE [dbo].[RepoTopics] DROP CONSTRAINT [FK_RepoTags_Project]

    -- RelatedCommunities
    ALTER TABLE [dbo].[RelatedCommunities] DROP CONSTRAINT [PK_RelatedCommunities]

    -- RepoOwners
    ALTER TABLE [dbo].[RepoOwners] DROP CONSTRAINT [PK_RepoOwner]
    ALTER TABLE [dbo].[RepoOwners] DROP CONSTRAINT [FK_RepoOwners_Projects] 
END

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


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

/* ActivityTypes > ActivityType */
IF (EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ActivityType') AND
    NOT EXISTS (SELECT * 
                 FROM INFORMATION_SCHEMA.TABLES 
                 WHERE TABLE_SCHEMA = 'dbo' 
                 AND  TABLE_NAME = 'ActivityTypes') )
BEGIN
    -- ApprovalRequestApprovers	ApprovalRequestApprover
    ALTER TABLE [dbo].[ApprovalRequestApprover] ADD CONSTRAINT [PK_ApprovalRequestApprover] PRIMARY KEY ([ApprovalRequestId], [ApproverEmail])
    ALTER TABLE [dbo].[ApprovalRequestApprover] ADD CONSTRAINT [FK_ApprovalRequestApprover_RepositoryApproval] FOREIGN KEY ([ApprovalRequestId]) REFERENCES [dbo].[RepositoryApproval]([Id])
    ALTER TABLE [dbo].[ApprovalRequestApprover] ADD CONSTRAINT [FK_ApprovalRequestApprover_User] FOREIGN KEY ([ApproverEmail]) REFERENCES [dbo].[User]([UserPrincipalName])

    -- ApprovalTypes	RepositoryApprovalType
    ALTER TABLE [dbo].[RepositoryApprovalType] ADD CONSTRAINT [FK_RepositoryApprovalType_User] FOREIGN KEY ([ApproverUserPrincipalName]) REFERENCES [dbo].[User]([UserPrincipalName])

    -- Approvers	RepositoryApprover
    ALTER TABLE [dbo].[RepositoryApprover] ADD CONSTRAINT [PK_RepositoryApprover] PRIMARY KEY ([ApprovalTypeId], [ApproverEmail])
    ALTER TABLE [dbo].[RepositoryApprover] ADD CONSTRAINT [FK_RepositoryApprover_RepositoryApprovalType] FOREIGN KEY ([ApprovalTypeId]) REFERENCES [dbo].[RepositoryApprovalType]([Id])
    ALTER TABLE [dbo].[RepositoryApprover] ADD CONSTRAINT [FK_RepositoryApprover_User] FOREIGN KEY ([ApproverEmail]) REFERENCES [dbo].[User]([UserPrincipalName])

    -- CategoryArticles	GuidanceCategoryArticle
    ALTER TABLE [dbo].[GuidanceCategoryArticle] ADD CONSTRAINT [FK_GuidanceCategoryArticle_GuidanceCategory] FOREIGN KEY([CategoryId]) REFERENCES [dbo].[GuidanceCategory]([Id])

    -- Communities	Community
    ALTER TABLE [dbo].[Community] ADD CONSTRAINT [FK_Community_ApprovalStatus] FOREIGN KEY ([ApprovalStatusId]) REFERENCES [dbo].[ApprovalStatus]([Id])

    -- CommunityActivities	CommunityActivity
    ALTER TABLE [dbo].[CommunityActivity] ADD CONSTRAINT [FK_CommunityActivity_Community] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Community]([Id])
    ALTER TABLE [dbo].[CommunityActivity] ADD CONSTRAINT [FK_CommunityActivity_ActivityType] FOREIGN KEY ([ActivityTypeId]) REFERENCES [dbo].[ActivityType]([Id])

    -- CommunityActivitiesContributionAreas	CommunityActivityContributionArea
    ALTER TABLE [dbo].[CommunityActivityContributionArea] ADD CONSTRAINT [FK_CommunityActivityContributionArea_CommunityActivity] FOREIGN KEY ([CommunityActivityId]) REFERENCES [dbo].[CommunityActivity]([Id])
    ALTER TABLE [dbo].[CommunityActivityContributionArea] ADD CONSTRAINT [FK_CommunityActivityContributionArea_ContributionArea] FOREIGN KEY ([ContributionAreaId]) REFERENCES [dbo].[ContributionArea]([Id])

    -- CommunityApprovalRequests	CommunityApprovalRequest
    ALTER TABLE [dbo].[CommunityApprovalRequest] ADD CONSTRAINT [PK_CommunityApprovalRequest] PRIMARY KEY ([CommunityId], [RequestId])
    ALTER TABLE [dbo].[CommunityApprovalRequest] ADD CONSTRAINT [FK_CommunityApprovalRequest_Community] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Community]([Id])
    ALTER TABLE [dbo].[CommunityApprovalRequest] ADD CONSTRAINT [FK_CommunityApprovalRequest_ApprovalRequest] FOREIGN KEY ([RequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])

    -- CommunityApprovals	ApprovalRequest
    ALTER TABLE [dbo].[ApprovalRequest] ADD CONSTRAINT [FK_ApprovalRequest_User] FOREIGN KEY ([ApproverUserPrincipalName]) REFERENCES [dbo].[User]([UserPrincipalName])
    ALTER TABLE [dbo].[ApprovalRequest] ADD CONSTRAINT [FK_ApprovalRequest_ApprovalStatus] FOREIGN KEY ([ApprovalStatusId]) REFERENCES [dbo].[ApprovalStatus]([Id])

    -- CommunityApproversList	CommunityApprover
    ALTER TABLE [dbo].[CommunityApprover] ADD CONSTRAINT [FK_CommunityApprover_User] FOREIGN KEY ([ApproverUserPrincipalName]) REFERENCES [dbo].[User]([UserPrincipalName])
    ALTER TABLE [dbo].[CommunityApprover] ADD CONSTRAINT [AK_ApproverUserPrincipalName_Categoy] UNIQUE ([ApproverUserPrincipalName], [Category])

    -- CommunityMembers	CommunityMember
    ALTER TABLE [dbo].[CommunityMember] ADD CONSTRAINT [PK_CommunityMember] PRIMARY KEY ([CommunityId], [UserPrincipalName])
    ALTER TABLE [dbo].[CommunityMember] ADD CONSTRAINT [FK_CommunityMember_Community] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Community]([Id])

    -- CommunitySponsors	CommunitySponsor
    ALTER TABLE [dbo].[CommunitySponsor] ADD CONSTRAINT [FK_CommunitySponsor_Community] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Community]([Id])
    ALTER TABLE [dbo].[CommunitySponsor] ADD CONSTRAINT [FK_CommunitySponsor_User] FOREIGN KEY ([UserPrincipalName]) REFERENCES [dbo].[User]([UserPrincipalName])
    ALTER TABLE [dbo].[CommunitySponsor] ADD CONSTRAINT [AK_CommunityId_UserPrincipalName] UNIQUE ([CommunityId], [UserPrincipalName])

    -- CommunityTags	CommunityTag
    ALTER TABLE [dbo].[CommunityTag] ADD CONSTRAINT [FK_CommunityTag_Community] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Community]([Id])
    ALTER TABLE [dbo].[CommunityTag] ADD CONSTRAINT [AK_CommunityId_Tag] UNIQUE ([CommunityId], [Tag])

    -- GitHubCopilot	GitHubCopilot
    ALTER TABLE [dbo].[GitHubCopilot] ADD CONSTRAINT [FK_GitHubCopilot_RegionalOrganization] FOREIGN KEY ([Region]) REFERENCES [dbo].[RegionalOrganization]([Id])

    -- GitHubCopilotApprovalRequests	GitHubCopilotApprovalRequest
    ALTER TABLE [dbo].[GitHubCopilotApprovalRequest] ADD CONSTRAINT [PK_GitHubCopilotApprovalRequest] PRIMARY KEY ([GitHubCopilotId], [RequestId])
    ALTER TABLE [dbo].[GitHubCopilotApprovalRequest] ADD CONSTRAINT [FK_GitHubCopilotApprovalRequest_GitHubCopilot] FOREIGN KEY ([GitHubCopilotId]) REFERENCES [dbo].[GitHubCopilot]([Id])
    ALTER TABLE [dbo].[GitHubCopilotApprovalRequest] ADD CONSTRAINT [FK_GitHubCopilotApprovalRequest_ApprovalRequest] FOREIGN KEY ([RequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])

    -- OrganizationAccess	OrganizationAccess
    ALTER TABLE [dbo].[OrganizationAccess] ADD CONSTRAINT [FK_OrganizationAccess_User] FOREIGN KEY ([UserPrincipalName]) REFERENCES [dbo].[User]([UserPrincipalName])
    ALTER TABLE [dbo].[OrganizationAccess] ADD CONSTRAINT [FK_OrganizationAccess_RegionalOrganization] FOREIGN KEY ([OrganizationId]) REFERENCES [dbo].[RegionalOrganization]([Id])

    -- OrganizationAccessApprovalRequests	OrganizationAccessApprovalRequest
    ALTER TABLE [dbo].[OrganizationAccessApprovalRequest] ADD CONSTRAINT [PK_OrganizationAccessApprovalRequest] PRIMARY KEY ([OrganizationAccessId], [RequestId])
    ALTER TABLE [dbo].[OrganizationAccessApprovalRequest] ADD CONSTRAINT [FK_OrganizationAccessApprovalRequest_OrganizationAccess] FOREIGN KEY ([OrganizationAccessId]) REFERENCES [dbo].[OrganizationAccess]([Id])
    ALTER TABLE [dbo].[OrganizationAccessApprovalRequest] ADD CONSTRAINT [FK_OrganizationAccessApprovalRequest_ApprovalRequest] FOREIGN KEY ([RequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])

    -- OrganizationApprovalRequests	OrganizationApprovalRequest
    ALTER TABLE [dbo].[OrganizationApprovalRequest] ADD CONSTRAINT [PK_OrganizationApprovalRequest] PRIMARY KEY ([OrganizationId], [RequestId])
    ALTER TABLE [dbo].[OrganizationApprovalRequest] ADD CONSTRAINT [FK_OrganizationApprovalRequest_Organization] FOREIGN KEY ([OrganizationId]) REFERENCES [dbo].[Organization]([Id])
    ALTER TABLE [dbo].[OrganizationApprovalRequest] ADD CONSTRAINT [FK_OrganizationApprovalRequest_ApprovalRequest] FOREIGN KEY ([RequestId]) REFERENCES [dbo].[ApprovalRequest]([Id])

    -- Organizations	Organization
    ALTER TABLE [dbo].[Organization] ADD CONSTRAINT [FK_Organization_RegionalOrganization] FOREIGN KEY ([Region]) REFERENCES [dbo].[RegionalOrganization]([Id])

    -- ProjectApprovals	RepositoryApproval
    ALTER TABLE [dbo].[RepositoryApproval] ADD CONSTRAINT [FK_RepositoryApproval_Repository] FOREIGN KEY ([ProjectId]) REFERENCES [dbo].[Repository]([Id])
    ALTER TABLE [dbo].[RepositoryApproval] ADD CONSTRAINT [FK_RepositoryApproval_RepositoryApprovalType] FOREIGN KEY ([ApprovalTypeId]) REFERENCES [dbo].[RepositoryApprovalType]([Id])
    ALTER TABLE [dbo].[RepositoryApproval] ADD CONSTRAINT [FK_RepositoryApproval_ApprovalStatus] FOREIGN KEY ([ApprovalStatusId]) REFERENCES [dbo].[ApprovalStatus]([Id])

    -- Projects	Repository
    ALTER TABLE [dbo].[Repository] ADD CONSTRAINT [FK_Repository_ApprovalStatus] FOREIGN KEY ([ApprovalStatusId]) REFERENCES [dbo].[ApprovalStatus]([Id])
    ALTER TABLE [dbo].[Repository] ADD CONSTRAINT [FK_Repository_Visibility] FOREIGN KEY ([VisibilityId]) REFERENCES [dbo].[Visibility]([Id])
    ALTER TABLE [dbo].[Repository] ADD CONSTRAINT [FK_Repository_OSSContributionSponsor] FOREIGN KEY ([OSSContributionSponsorId]) REFERENCES [dbo].[OSSContributionSponsor]([Id])

    -- RelatedCommunities	RelatedCommunity
    ALTER TABLE [dbo].[RelatedCommunity] ADD CONSTRAINT [PK_RelatedCommunity] PRIMARY KEY ([ParentCommunityId], [RelatedCommunityId])

    -- RepoOwners	RepositoryOwner
    ALTER TABLE [dbo].[RepositoryOwner] ADD CONSTRAINT [PK_RepositoryOwner] PRIMARY KEY ([ProjectId], [UserPrincipalName])
    ALTER TABLE [dbo].[RepositoryOwner] ADD CONSTRAINT [FK_RepositoryOwner_Repository] FOREIGN KEY ([ProjectId]) REFERENCES [dbo].[Repository]([Id])

    -- RepoTopics	RepositoryTopic
    ALTER TABLE [dbo].[RepositoryTopic] ADD CONSTRAINT [PK_RepositoryTopic] PRIMARY KEY ([Topic], [ProjectId])
    ALTER TABLE [dbo].[RepositoryTopic] ADD CONSTRAINT [FK_RepositoryTopic_Repository] FOREIGN KEY ([ProjectId]) REFERENCES [dbo].[Repository]([Id])
END
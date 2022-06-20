SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[PR_CommunitiesActivities_Select]
AS
BEGIN
    SET NOCOUNT ON
    SELECT TOP (1000) [c].[Id]
      ,[ca].[Name]
      ,[CommunityId]
      ,[c].[Name] AS 'CommunityName'
      ,[ActivityTypeId]
      ,[a].[Name] AS 'TypeName'
      ,[car].[Id] AS 'PrimaryContributionAreaId'
      ,[car].[Name] AS 'PrimaryContributionAreaName'
      ,[c].[Url]
      ,[Date]
      ,[c].[Created]
      ,[c].[CreatedBy]
      ,[c].[Modified]
      ,[c].[ModifiedBy]
    FROM [dbo].[CommunityActivities] AS ca
    INNER JOIN [dbo].[Communities] AS c ON ca.CommunityId = c.Id
    INNER JOIN [dbo].[ActivityTypes] AS a ON a.Id = c.Id
    INNER JOIN (
      SELECT * FROM [dbo].[CommunityActivitiesContributionAreas] WHERE IsPrimary = 1
    ) AS caca ON caca.CommunityActivityId = ca.Id
    INNER JOIN [dbo].[ContributionAreas] AS car ON car.Id = caca.ContributionAreaId
END
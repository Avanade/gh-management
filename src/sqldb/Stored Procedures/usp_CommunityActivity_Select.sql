CREATE PROCEDURE [dbo].[usp_CommunityActivity_Select]
AS
BEGIN
  SET NOCOUNT ON
  SELECT
    [CA].[Id],
    [CA].[Name],
    [CommunityId],
    [C].[Name] AS [CommunityName],
    [ActivityTypeId],
    [A].[Name] AS [ActivityTypeName],
    [CAR].[Id] AS [PrimaryContributionAreaId],
    [CAR].[Name] AS [PrimaryContributionAreaName],
    [C].[Url],
    [Date],
    [CA].[Created],
    [CA].[CreatedBy],
    [CA].[Modified],
    [CA].[ModifiedBy]
  FROM [dbo].[CommunityActivity] AS [CA]
    LEFT JOIN [dbo].[Community] AS [C] ON [CA].[CommunityId] = [C].[Id]
    LEFT JOIN [dbo].[ActivityType] AS [A] ON [CA].[ActivityTypeId] = [A].[Id]
    LEFT JOIN (
    SELECT *
    FROM [dbo].[CommunityActivityContributionArea]
    WHERE [IsPrimary] = 1
  ) AS [CACA] ON [CACA].[CommunityActivityId] = [CA].[Id]
    LEFT JOIN [dbo].[ContributionArea] AS [CAR] ON [CAR].[Id] = [CACA].[ContributionAreaId]
  ORDER BY [CA].[Modified] ASC
END
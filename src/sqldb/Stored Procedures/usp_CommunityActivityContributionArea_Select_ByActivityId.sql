CREATE PROCEDURE [dbo].[usp_CommunityActivityContributionArea_Select_ByActivityId]
	@ActivityId [INT]
AS
BEGIN
  SET NOCOUNT ON
  SELECT 
    [CACA].[Id],
    [CACA].[CommunityActivityId],
    [CACA].[ContributionAreaId],
    [CA].[Name] AS [ContributionAreaName],
    [CACA].[IsPrimary]
  FROM 
    [dbo].[CommunityActivityContributionArea] AS [CACA] 
  JOIN
    [dbo].[ContributionArea] AS [CA] ON [CACA].[ContributionAreaId] = [CA].[Id]
  WHERE
    [CACA].[CommunityActivityId] = @ActivityId
END
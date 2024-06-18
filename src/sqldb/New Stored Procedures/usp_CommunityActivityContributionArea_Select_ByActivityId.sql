CREATE PROCEDURE [dbo].[usp_CommunityActivityContributionArea_Select_ByActivityId]
	@ActivityId [INT]
AS
BEGIN
  SET NOCOUNT ON
  SELECT 
    * 
  FROM 
    [dbo].[CommunityActivityContributionArea] AS [CACA] 
  JOIN
    [dbo].[ContributionArea] AS [CA] ON [CACA].[ContributionAreaId] = [CA].[Id]
  WHERE
    [CACA].[CommunityActivityId] = @ActivityId AND [CACA].[IsPrimary] = 0
END
CREATE PROCEDURE [dbo].[usp_CommunityActivity_TotalCount_ByOptionAndCreatedBy]
	@Search [VARCHAR](50) = '',
	@CreatedBy [VARCHAR](50)
AS
BEGIN
  SET NOCOUNT ON
	SELECT 
		COUNT(*) AS [Total]
  FROM [dbo].[CommunityActivity] AS [CA]
  LEFT JOIN [dbo].[Community] AS [C] ON [CA].[CommunityId] = [C].[Id]
  LEFT JOIN [dbo].[ActivityType] AS [A] ON [CA].[ActivityTypeId] = [A].[Id]
  LEFT JOIN (
    SELECT * FROM [dbo].[CommunityActivityContributionArea] WHERE [IsPrimary] = 1
  ) AS [CACA] ON [CACA].[CommunityActivityId] = [CA].[Id]
  LEFT JOIN [dbo].[ContributionArea] AS [CR] ON [CR].[Id] = [CACA].[ContributionAreaId]
  WHERE
  (
    [CA].[Name] LIKE '%'+@Search+'%' OR
    [C].[Name] LIKE '%'+@Search+'%' OR
    [A].[Name] LIKE '%'+@Search+'%' OR
    [CR].[Name] LIKE '%'+@Search+'%'
  ) AND [CA].[CreatedBy] = @CreatedBy
END
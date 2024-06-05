CREATE PROCEDURE [dbo].[usp_CommunityActivity_Select_ByOptionAndCreatedBy]
	@Offset [INT] = 0,
	@Filter [INT] = 10,
	@Search [VARCHAR](50) = '',
	@OrderBy [VARCHAR](50) = 'Date',
	@OrderType [VARCHAR](5) = 'ASC',
	@CreatedBy [VARCHAR](50)
AS
BEGIN
  SET NOCOUNT ON
	SELECT 
		[CA].[Id],
		[CA].[Name],
		[CommunityId],
		[C].[Name] AS [CommunityName],
		[ActivityTypeId],
		[A].[Name] AS [TypeName],
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
		SELECT * FROM [dbo].[CommunityActivityContributionArea] WHERE [IsPrimary] = 1
	) AS [CACA] ON [CACA].[CommunityActivityId] = [CA].[Id]
	LEFT JOIN [dbo].[ContributionArea] AS [CAR] ON [CAR].[Id] = [CACA].[ContributionAreaId]
	WHERE
	(
		[CA].[Name] LIKE '%' + @Search + '%' OR
		[C].[Name] LIKE '%' + @Search + '%' OR
		[A].[Name] LIKE '%' + @Search + '%' OR
		[CAR].[Name] LIKE '%' + @Search + '%'
	) AND [CA].[CreatedBy] = @CreatedBy
	ORDER BY
		CASE WHEN @OrderType='ASC' THEN
			CASE @OrderBy
				WHEN 'Date' THEN [CA].[Date]
			END
		END,
		CASE WHEN @OrderType='DESC' THEN
			CASE @OrderBy
				WHEN 'Date' THEN [CA].[Date]
			END
		END DESC,
		CASE WHEN @OrderType='ASC' THEN
			CASE @OrderBy
				WHEN 'Activity' THEN [CA].[Name]
				WHEN 'Community' THEN [C].[Name]
				WHEN 'Type' THEN [A].[Name]
				WHEN 'PrimaryContributionArea' THEN [CAR].[Name]
			END
		END,
		CASE WHEN @OrderType='DESC' THEN
			CASE @OrderBy
				WHEN 'Activity' THEN [CA].[Name]
				WHEN 'Community' THEN [C].[Name]
				WHEN 'Type' THEN [A].[Name]
				WHEN 'PrimaryContributionArea' THEN [CAR].[Name]
			END
		END DESC
	OFFSET @Offset ROWS 
	FETCH NEXT @Filter ROWS ONLY
END
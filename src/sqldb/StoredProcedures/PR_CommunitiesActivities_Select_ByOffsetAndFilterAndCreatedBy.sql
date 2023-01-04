CREATE PROCEDURE [dbo].[PR_CommunityActivities_Select_ByOffsetAndFilterAndCreatedBy](
	@Offset int = 0,
	@Filter int = 10,
	@Search varchar(50) = '',
	@OrderBy varchar(50) = 'Date',
	@OrderType varchar(5) = 'ASC',
	@CreatedBy varchar(50)
)
AS
BEGIN
    SET NOCOUNT ON
	SELECT [ca].[Id]
	  ,[ca].[Name]
      ,[CommunityId]
	  ,[c].[Name] AS 'CommunityName'
      ,[ActivityTypeId]
	  ,[a].[Name] AS 'TypeName'
	  ,[car].[Id] AS 'PrimaryContributionAreaId'
	  ,[car].[Name] AS 'PrimaryContributionAreaName'
      ,[c].[Url]
      ,[Date]
      ,[ca].[Created]
      ,[ca].[CreatedBy]
      ,[ca].[Modified]
      ,[ca].[ModifiedBy]
	  FROM [dbo].[CommunityActivities] AS ca
	  LEFT JOIN [dbo].[Communities] AS c ON ca.CommunityId = c.Id
	  LEFT JOIN [dbo].[ActivityTypes] AS a ON ca.ActivityTypeId = a.Id
	  LEFT JOIN (
		SELECT * FROM [dbo].[CommunityActivitiesContributionAreas] WHERE IsPrimary = 1
	  ) AS caca ON caca.CommunityActivityId = ca.Id
	  LEFT JOIN [dbo].[ContributionAreas] AS car ON car.Id = caca.ContributionAreaId
	  WHERE
		(
			ca.Name LIKE '%'+@Search+'%' OR
			c.Name LIKE '%'+@Search+'%' OR
			a.Name LIKE '%'+@Search+'%' OR
			car.Name LIKE '%'+@Search+'%'
		) AND ca.CreatedBy = @CreatedBy
	  ORDER BY
		CASE WHEN @OrderType='ASC' THEN
			CASE @OrderBy
				WHEN 'Date' THEN [ca].[Date]
			END
		END,
		CASE WHEN @OrderType='DESC' THEN
			CASE @OrderBy
				WHEN 'Date' THEN [ca].[Date]
			END
		END DESC,
		CASE WHEN @OrderType='ASC' THEN
			CASE @OrderBy
				WHEN 'Activity' THEN [ca].[Name]
				WHEN 'Community' THEN [c].[Name]
				WHEN 'Type' THEN [a].[Name]
				WHEN 'PrimaryContributionArea' THEN [car].[Name]
			END
		END,
		CASE WHEN @OrderType='DESC' THEN
			CASE @OrderBy
				WHEN 'Activity' THEN [ca].[Name]
				WHEN 'Community' THEN [c].[Name]
				WHEN 'Type' THEN [a].[Name]
				WHEN 'PrimaryContributionArea' THEN [car].[Name]
			END
		END DESC
	  OFFSET @Offset ROWS 
	  FETCH NEXT @Filter ROWS ONLY
END
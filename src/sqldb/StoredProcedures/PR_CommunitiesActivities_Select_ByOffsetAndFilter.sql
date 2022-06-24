/****** Object:  StoredProcedure [dbo].[PR_CommunityActivities_Select_ByOffsetAndFilter]    Script Date: 6/24/2022 11:37:16 AM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[PR_CommunityActivities_Select_ByOffsetAndFilter](
	@Offset int = 0,
	@Filter int = 10,
	@Search varchar(50) = '',
	@OrderBy varchar(5) = 'ASC'
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
		ca.Name LIKE '%'+@search+'%' OR
		c.Name LIKE '%'+@search+'%' OR
		a.Name LIKE '%'+@search+'%' OR
		car.Name LIKE '%'+@search+'%'
	  ORDER by 
		CASE WHEN @OrderBy='ASC' THEN ca.Modified  END,
		CASE WHEN @OrderBy='DESC' THEN ca.Modified  END DESC
	  OFFSET @Offset ROWS 
	  FETCH NEXT @Filter ROWS ONLY
END
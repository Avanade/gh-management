CREATE PROCEDURE [dbo].[PR_CommunityActivities_TotalCount_ByCreatedBy] (
	@Search varchar(50) = '',
	@CreatedBy varchar(50)
)
AS
BEGIN
    SET NOCOUNT ON
	SELECT 
		COUNT(*) AS Total 
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
END
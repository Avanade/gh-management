CREATE PROCEDURE [dbo].[PR_Repositories_Select_ByOffsetAndFilter](
	@Offset INT = 0,
	@Search VARCHAR(50) = ''
)
AS
BEGIN
    SET NOCOUNT ON
SELECT DISTINCT [p].[Id],
        [p].[Name],
        [P].[AssetCode],
        [p].[Organization],
        [p].[Description],
        [p].[IsArchived],
        [p].[Created],
        [p].[RepositorySource],
        [p].[TFSProjectReference],
        [v].[Name] AS "Visibility",
		    [p].[CoOwner],
        [p].[Createdby],
        (SELECT STRING_AGG(r.Topic, ',') FROM dbo.RepoTopics AS r WHERE r.ProjectId=p.Id) AS "Topics"
	  FROM [dbo].[Projects] AS p
	  LEFT JOIN [dbo].[Visibility] AS v ON p.VisibilityId = v.Id
    LEFT JOIN [dbo].[RepoTopics] AS rt ON rt.ProjectId = p.Id
    WHERE
		p.Name LIKE '%'+@search+'%' OR
		rt.Topic LIKE '%'+@search+'%'
    ORDER BY
		[p].[Created] DESC
	  OFFSET @Offset ROWS 
	  FETCH NEXT 15 ROWS ONLY
END
CREATE PROCEDURE [dbo].[PR_Repositories_Select_ByOffsetAndFilter](
	@Offset INT = 0,
	@Search VARCHAR(50) = ''
)
AS
BEGIN
  SET NOCOUNT ON
  SELECT  [p].[Id],
          [p].[Name],
          [P].[AssetCode],
          [p].[ECATTID],
          [p].[Organization],
          [p].[Description],
          [p].[IsArchived],
          [p].[Created],
          [p].[RepositorySource],
          [p].[TFSProjectReference],
          [v].[Name] AS "Visibility",
          [p].[CoOwner],
          [p].[Createdby],
          (SELECT STRING_AGG(r.Topic, ',') FROM dbo.RepoTopics AS r WHERE r.ProjectId=p.Id) AS "Topics",
          COUNT(*) AS Score
      FROM [dbo].[Projects] AS p
      LEFT JOIN [dbo].[Visibility] AS v ON p.VisibilityId = v.Id
      LEFT JOIN [dbo].[RepoTopics] AS rt ON rt.ProjectId = p.Id
      INNER JOIN STRING_SPLIT(@Search, ' ') AS ss ON (p.Name LIKE '%'+ss.[value]+'%' OR rt.Topic LIKE '%'+ss.[value]+'%')
      GROUP BY [p].[Id],
          [p].[Name],
          [P].[AssetCode],
          [p].[ECATTID],
          [p].[Organization],
          [p].[Description],
          [p].[IsArchived],
          [p].[Created],
          [p].[RepositorySource],
          [p].[TFSProjectReference],
          [v].[Name],
          [p].[CoOwner],
          [p].[Createdby]
      ORDER BY
      Score DESC,
      [p].[Created] DESC
    OFFSET @Offset ROWS 
    FETCH NEXT 15 ROWS ONLY
END
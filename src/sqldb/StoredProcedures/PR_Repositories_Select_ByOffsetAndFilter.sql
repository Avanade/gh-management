CREATE PROCEDURE [dbo].[PR_Repositories_Select_ByOffsetAndFilter](
	@Offset INT = 0,
	@Search VARCHAR(50) = ''
)
AS
BEGIN
    SET NOCOUNT ON
	SELECT [p].[Id],
        [p].[Name],
        [p].[Description],
        [p].[IsArchived],
        [p].[Created],
        [p].[RepositorySource],
        [p].[TFSProjectReference],
        [v].[Name] AS "Visibility"
	  FROM [dbo].[Projects] AS p
	  LEFT JOIN [dbo].[Visibility] AS v ON p.VisibilityId = v.Id
    WHERE
		p.Name LIKE '%'+@search+'%'
    ORDER BY
		[p].[Name]
	  OFFSET @Offset ROWS 
	  FETCH NEXT 15 ROWS ONLY
END
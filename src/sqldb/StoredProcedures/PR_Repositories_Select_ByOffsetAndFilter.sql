SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[PR_Repositories_Select_ByOffsetAndFilter](
	@Offset int = 0,
	@Search varchar(50) = ''
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
        [v].[Name] as "Visibility"
	  FROM [dbo].[Projects] AS p
	  LEFT JOIN [dbo].[Visibility] AS v ON p.VisibilityId = v.Id
    WHERE
		p.Name LIKE '%'+@search+'%'
    ORDER BY
		[p].[Name]
	  OFFSET @Offset ROWS 
	  FETCH NEXT 15 ROWS ONLY
END
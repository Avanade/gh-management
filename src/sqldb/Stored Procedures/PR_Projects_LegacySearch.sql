CREATE PROCEDURE [dbo].[PR_Projects_LegacySearch](
	@searchText VARCHAR (100)
)
AS 
	SELECT
		[AssetCode] AS 'Code',
		[Name] AS 'Title',
		[Description],
		'Asset' [Type]
	FROM 
		[dbo].[Projects]
	WHERE	
		[Name] LIKE '%'+@searchText+'%' AND
        [RepositorySource] != 'GitHub'

UNION

	SELECT
		CAST(GitHubId as varchar(50)) AS 'Code',
		[Name] AS 'Title',
		[Description],
		'Asset' [Type]
	FROM 
		[dbo].[Projects]
	WHERE	
		[Name] LIKE '%'+@searchText+'%' AND
        [RepositorySource] = 'GitHub'

ORDER BY [Name]
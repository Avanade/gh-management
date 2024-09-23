CREATE PROCEDURE [dbo].[usp_Repository_LegacySearch]
	@Search [VARCHAR] (100)
AS
BEGIN 
	SELECT
		[AssetCode] AS [Code],
		[Name] AS [Title],
		[Description],
		'Asset' AS [Type]
	FROM 
		[dbo].[Repository]
	WHERE	
		[Name] LIKE '%' + @Search + '%' AND
    [RepositorySource] != 'GitHub' AND
    [VisibilityId] != 1
  
  UNION

  SELECT
    CAST([GithubId] AS VARCHAR(50)) AS [Code],
    [Name] AS [Title],
    [Description],
    'Asset' AS [Type]
  FROM 
    [dbo].[Repository]
  WHERE	
    [Name] LIKE '%' + @Search + '%' AND
    [RepositorySource] = 'GitHub' AND
    [VisibilityId] != 1

  ORDER BY [Name]
END
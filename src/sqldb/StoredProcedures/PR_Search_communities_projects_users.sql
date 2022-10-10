CREATE PROCEDURE [dbo].[PR_Search_communities_projects_users]

@searchText VARCHAR (100),
@offSet INT = 0,
@rowCount INT = 0

AS 

SELECT
	'Users' [Source],[Name],
	 UserPrincipalName [Description],
	 Users.GitHubId [Id]
FROM [dbo].[Users]
WHERE [Name] LIKE '%'+@searchText+'%'
OR [UserPrincipalName] LIKE '%'+@searchText+'%'


UNION

SELECT 
	'Projects' [Source], [Name],
	CASE
		WHEN [CreatedBy] IS NULL THEN [RepositorySource]
		ELSE [RepositorySource] + ' - ' + [CreatedBy]
	END [Description],
	Projects.Id [ID]
FROM [dbo].[Projects]
WHERE [Name] LIKE '%'+@searchText+'%'
OR [CreatedBy] LIKE '%'+@searchText+'%'

UNION

SELECT 
	'Communities' [Source],[Name],
	[Description],
	Communities.[Id]
FROM [dbo].[Communities]
WHERE [Name] LIKE '%'+@searchText+'%'
OR [Description] LIKE '%'+@searchText+'%'

ORDER BY [Name]
OFFSET @offSet ROWS
FETCH NEXT @rowCount ROWS ONLY
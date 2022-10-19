CREATE PROCEDURE [dbo].[PR_Search_communities_projects_users]

@searchText VARCHAR (100),
@offSet INT = 0,
@rowCount INT = 0

AS 

SELECT
		'Users' [Source],[Name],
		UserPrincipalName [Description],
		Users.GitHubId [Id]
FROM	[dbo].[Users]
WHERE	[Name] LIKE '%'+@searchText+'%'
		OR [UserPrincipalName] LIKE '%'+@searchText+'%'

UNION

SELECT 
		'Repositories' [Source], [Name],
			CASE
				WHEN [CreatedBy] IS NULL THEN [RepositorySource]
				ELSE [RepositorySource] + ' - ' + [CreatedBy]
			END [Description],
				Projects.Id [ID]
FROM	[dbo].[Projects]
WHERE	[Name] LIKE '%'+@searchText+'%'
		OR [CreatedBy]
		LIKE '%'+@searchText+'%'

UNION

SELECT 
		'Communities' [Source],c.[Name],
		[Description],
		c.[Id]
FROM	[dbo].[Communities] c
  	INNER JOIN ApprovalStatus T ON c.ApprovalStatusId = T.Id
WHERE	(c.[Name] LIKE '%'+@searchText+'%'
		OR [Description] LIKE '%'+@searchText+'%' )

		AND 
		(c.ApprovalStatusId = 5 OR c.CreatedBy = 'dennis.delamida@accenture.com')
				AND
		c.ApprovalStatusId =5

ORDER BY [Name]
OFFSET @offSet ROWS
FETCH NEXT @rowCount ROWS ONLY
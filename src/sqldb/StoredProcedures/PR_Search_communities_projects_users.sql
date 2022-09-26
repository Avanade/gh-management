SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

Create Procedure [dbo].[PR_Search_communities_projects_users]

@searchText varchar (100),
@offSet int = 0,
@rowCount int = 0

as 

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
	case
		when [CreatedBy] IS NULL then [RepositorySource]
		else [RepositorySource] + ' - ' + [CreatedBy]
	end [Description],
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
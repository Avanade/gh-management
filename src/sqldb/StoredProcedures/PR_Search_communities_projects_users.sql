SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

ALTER Procedure [dbo].[PR_Search_communities_projects_users]

@searchText varchar (100),
@offSet int = 0,
@rowCount int = 0

as 

SELECT
	'Users' [Source],[Name],
	 UserPrincipalName [Description]
FROM [dbo].[Users]
WHERE [Name] LIKE '%'+@searchText+'%'
OR [UserPrincipalName] LIKE '%'+@searchText+'%'


UNION

SELECT 
	'Projects' [Source], [Name],
	CoOwner [Description]
FROM [dbo].[Projects]
WHERE [Name] LIKE '%'+@searchText+'%'
OR [CoOwner] LIKE '%'+@searchText+'%'

UNION

SELECT 
	'Communities' [Source],[Name],
	 [Description]
FROM [dbo].[Communities]
WHERE [Name] LIKE '%'+@searchText+'%'
OR [Description] LIKE '%'+@searchText+'%'

ORDER BY [Name]
OFFSET @offSet ROWS
FETCH NEXT @rowCount ROWS ONLY
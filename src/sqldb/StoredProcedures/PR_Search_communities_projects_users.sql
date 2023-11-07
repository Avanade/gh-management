CREATE PROCEDURE [dbo].[PR_Search_communities_projects_users]
	@searchText VARCHAR (100),
	@offSet INT = 0,
	@rowCount INT = 0,
	@userprincipal VARCHAR (100) = null
AS 
	SELECT
		'Users' [Source],
		[Name],
		CONCAT(
			'User Principal Name: ', [UserPrincipalName], ',',
			'Github ID: ', CASE WHEN [GitHubId] IS NULL THEN 'N/A' ELSE [GitHubId] END, ',',
			'Github User: ', CASE WHEN [GitHubUser] IS NULL THEN 'N/A' ELSE [GitHubUser] END
		) [Description],
		Users.GitHubId [Id]
	FROM 
		[dbo].[Users]
	WHERE	
		[Name] LIKE '%'+@searchText+'%' OR 
		[UserPrincipalName] LIKE '%'+@searchText+'%' OR
		[GitHubId] LIKE '%'+@searchText+'%' OR
		[GitHubUser] LIKE '%'+@searchText+'%'
UNION
	SELECT 
		'Repositories' [Source], 
		[Name],
		CONCAT(
			CASE
			WHEN [CreatedBy] IS NULL THEN 
				[RepositorySource]
			ELSE 
				[RepositorySource] + ' - ' + [CreatedBy]
			END, '|', 
			(
				SELECT 
					STRING_AGG(Topic, ',') 
				FROM 
					RepoTopics 
				WHERE ProjectId=Id
			)
		) [Description],
		Id [ID]
	FROM (
		SELECT 
			Name, 
			RepositorySource, 
			CreatedBy, 
			Id
		FROM 
			Projects AS P 
		LEFT JOIN 
			RepoOwners AS RO ON RO.ProjectId = P.Id
		LEFT JOIN 
			RepoTopics RT ON RT.ProjectId = P.Id
		WHERE	
			[Name] LIKE '%'+@searchText+'%' OR 
			RO.UserPrincipalName LIKE '%'+@searchText+'%' OR 
			RT.Topic LIKE '%'+@searchText+'%'
	) AS Repository
	GROUP BY Name, CreatedBy, RepositorySource, Id
UNION
	SELECT 
		'Communities' [Source],
		c.[Name],
		[Description],
		c.[Id]
	FROM	
		[dbo].[Communities] c
	LEFT JOIN 
		ApprovalStatus T ON c.ApprovalStatusId = T.Id
	WHERE 
		(
			(
				c.[Name] LIKE '%'+@searchText+'%' OR 
				[Description] LIKE '%'+@searchText+'%' 
			) AND 
			c.ApprovalStatusId = 5
		) OR
		(
			(
				c.[Name] LIKE '%'+@searchText+'%' OR 
				[Description] LIKE '%'+@searchText+'%' 
			) AND 
			c.CreatedBy = @userprincipal
		)
ORDER BY [Name]
OFFSET @offSet ROWS
FETCH NEXT @rowCount ROWS ONLY
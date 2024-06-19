CREATE PROCEDURE [dbo].[PR_Search_communities_projects_users]
	@searchText VARCHAR (100),
	@offSet INT = 0,
	@rowCount INT = 0,
	@userprincipal VARCHAR (100) = null
AS 
	(SELECT
		'Users' [Source],
		[Name],
		CONCAT(
			'User Principal Name: ', [UserPrincipalName], ',',
			'Github ID: ', CASE WHEN [GitHubId] IS NULL THEN 'N/A' ELSE [GitHubId] END, ',',
			'Github User: ', CASE WHEN [GitHubUser] IS NULL THEN 'N/A' ELSE [GitHubUser] END
		) [Description],
		[GitHubId] [Id],
		COUNT(*) [Score]
	FROM 
		[dbo].[Users] AS U
		JOIN STRING_SPLIT(@searchText, ' ') AS SS ON (
			[Name] LIKE '%'+ss.Value+'%' OR 
			[UserPrincipalName] LIKE '%'+ss.Value+'%' OR
			[GitHubId] LIKE '%'+ss.Value+'%' OR
			[GitHubUser] LIKE '%'+ss.Value+'%'
		)
	GROUP BY [GitHubId], [Name], [UserPrincipalName], [GitHubUser])
UNION
	(SELECT 
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
		Id [ID],
		COUNT(*) Score
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
		JOIN STRING_SPLIT(@searchText, ' ') AS SS ON (	
			[Name] LIKE '%'+SS.value+'%' OR 
			RO.UserPrincipalName LIKE '%'+SS.value+'%' OR 
			RT.Topic LIKE '%'+SS.value+'%'
		)
	) AS Repository
	GROUP BY Name, CreatedBy, RepositorySource, Id)
UNION
	(SELECT 
		'Communities' [Source],
		c.[Name],
		[Description],
		c.[Id],
		COUNT(*) Score
	FROM	
		[dbo].[Communities] c
	LEFT JOIN 
		ApprovalStatus T ON c.ApprovalStatusId = T.Id
	JOIN STRING_SPLIT(@searchText, ' ') AS SS ON (
		c.[Name] LIKE '%'+SS.value+'%' OR 
		[Description] LIKE '%'+SS.value+'%'
	)
	WHERE
		c.ApprovalStatusId = 5 AND c.CreatedBy = @userprincipal
	GROUP BY c.Name, c.Description, c.Id)
ORDER BY Score DESC
OFFSET @offSet ROWS
FETCH NEXT @rowCount ROWS ONLY
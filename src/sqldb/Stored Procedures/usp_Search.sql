CREATE PROCEDURE [dbo].[usp_Search]
  @Search [VARCHAR] (100),
  @OffSet [INT] = 0,
  @RowCount [INT] = 0,
  @UserPrincipalName [VARCHAR] (100) = null
AS
BEGIN
  (
    SELECT
      'Users' AS [Source],
      [Name],
      CONCAT(
        'User Principal Name: ', [UserPrincipalName], ',',
        'Github ID: ', CASE WHEN [GitHubId] IS NULL THEN 'N/A' ELSE [GitHubId] END, ',',
        'Github User: ', CASE WHEN [GitHubUser] IS NULL THEN 'N/A' ELSE [GitHubUser] END
      ) AS [Description],
      [GitHubId] AS [Id],
      COUNT(*) AS [Score]
    FROM [dbo].[User] AS [U]
    JOIN STRING_SPLIT(@Search, ' ') AS [SS] ON (
      [Name] LIKE '%' + [SS].[value] + '%' OR
      [UserPrincipalName] LIKE '%' + [SS].[value] + '%' OR
      [GitHubId] LIKE '%' + [SS].[value] + '%' OR
      [GitHubUser] LIKE '%' + [SS].[value] + '%'
    )
    GROUP BY [GitHubId], [Name], [UserPrincipalName], [GitHubUser]
  )
  UNION
  (
    SELECT
      'Repositories' AS [Source],
      CONCAT (
        [Organization], '/',
        [Name]
      ),
      CONCAT (
        [Description], '|', (
          SELECT STRING_AGG([Topic], ',')
          FROM [dbo].[RepositoryTopic]
          WHERE [RepositoryId] = [Id]
        )
      ) AS [Description],
      [Id] AS [ID],
      COUNT(*) AS [Score]
    FROM (
      SELECT
        [Name],
        [Description],
        [Organization],
        [Id]
      FROM [dbo].[Repository] AS [P]
      LEFT JOIN [dbo].[RepositoryTopic] AS [RT] ON [RT].[RepositoryId] = [P].[Id]
      JOIN STRING_SPLIT(@Search, ' ') AS [SS] ON (	
        [Name] LIKE '%' + [SS].[value] + '%' OR
        [RT].[Topic] LIKE '%' + [SS].[value] + '%'
      ) 
      WHERE
        [P].[VisibilityId] != 1
    ) AS [Repository]
    GROUP BY [Name], [Description], [Organization], [Id]
  )
  UNION
  (
    SELECT
      'Communities' AS [Source],
      [C].[Name],
      [Description],
      [C].[Id],
      COUNT(*) AS [Score]
    FROM
      [dbo].[Community] AS [C]
      LEFT JOIN [dbo].[ApprovalStatus] AS [T] ON [C].[ApprovalStatusId] = [T].[Id]
      JOIN STRING_SPLIT(@Search, ' ') AS [SS] ON (
        [C].[Name] LIKE '%' + [SS].[value] + '%' OR
        [Description] LIKE '%' + [SS].[value] + '%'
      )
    WHERE [C].[ApprovalStatusId] = 5 AND [C].[CreatedBy] = @UserPrincipalName
    GROUP BY [C].[Name], [C].[Description], [C].[Id]
  )
  ORDER BY [Score] DESC
  OFFSET @OffSet ROWS
  FETCH NEXT @RowCount ROWS ONLY
END
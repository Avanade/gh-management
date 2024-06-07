CREATE PROCEDURE [dbo].[usp_RepositoryOwner_Select]
AS
BEGIN
  SELECT
    [RO].[RepositoryId],
    [P].[Name],
    [RO].[UserPrincipalName]
  FROM [dbo].[RepositoryOwner] AS [RO]
  LEFT JOIN [dbo].[Repository] AS [P] on [RO].[RepositoryId] = [P].[Id]
END
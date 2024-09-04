CREATE PROCEDURE [dbo].[usp_RepositoryOwner_Select_ByRepositoryId]
  @RepositoryId [INT]
AS
BEGIN
  SELECT
    [RO].[RepositoryId],
    [RO].[UserPrincipalName],
    [U].[GitHubUser]
  FROM [dbo].[RepositoryOwner] AS [RO]
    LEFT JOIN [dbo].[User]  AS [U] ON [RO].[UserPrincipalName] = [U].[UserPrincipalName]
  WHERE [RepositoryId] = @RepositoryId
END
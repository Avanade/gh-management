CREATE PROCEDURE [dbo].[usp_User_Select_ByGitHubId]
  @GitHubId [VARCHAR](100)
AS
BEGIN
  SET NOCOUNT ON;

  SELECT
    [UserPrincipalName],
    [Name],
    [GivenName],
    [SurName],
    [JobTitle],
    [GitHubId],
    [GitHubUser],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy],
    [LastGithubLogin]
  FROM [dbo].[User]
  WHERE [GitHubId] = @GitHubId
END
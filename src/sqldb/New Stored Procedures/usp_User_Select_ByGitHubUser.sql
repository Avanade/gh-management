CREATE PROCEDURE [dbo].[usp_User_Select_ByGitHubUser]
	@GitHubUser VARCHAR(100)
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
  WHERE [GitHubUser] = @GitHubUser
END
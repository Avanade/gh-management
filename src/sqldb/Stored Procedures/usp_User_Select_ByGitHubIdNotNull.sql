CREATE PROCEDURE [dbo].[usp_User_Select_ByGitHubIdNotNull]
AS
BEGIN
  SELECT
    [UserPrincipalName],
    [GivenName],
    [SurName],
    [JobTitle],
    [GitHubId],
    [GitHubUser],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy]
  FROM [dbo].[User]
  WHERE [GitHubId] IS NOT NULL
END
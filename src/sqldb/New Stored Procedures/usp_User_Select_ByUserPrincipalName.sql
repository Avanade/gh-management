CREATE PROCEDURE [dbo].[usp_User_Select_ByUserPrincipalName]
	@UserPrincipalName VARCHAR(100)
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
    [ModifiedBy]
  FROM [dbo].[User]
  WHERE [UserPrincipalName] = @UserPrincipalName
END
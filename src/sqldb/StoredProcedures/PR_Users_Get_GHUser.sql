CREATE PROCEDURE PR_Users_Get_GHUser

	@UserPrincipalName VARCHAR(100)

AS

SELECT GitHubId, GitHubUser FROM Users WHERE UserPrincipalName = @UserPrincipalName
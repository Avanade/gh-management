create PROCEDURE [dbo].[PR_Users_GetEmailByGitHubUsername]
(
	@GithubUser varchar(100)
)
AS
BEGIN
SELECT 
		[UserPrincipalName]
  FROM 
		[dbo].[Users]
  WHERE  
	    [GithubUser] = @GithubUser

END
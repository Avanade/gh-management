create PROCEDURE [dbo].[PR_Users_GetEmail]
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
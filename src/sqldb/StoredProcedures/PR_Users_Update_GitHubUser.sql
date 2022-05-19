CREATE PROCEDURE [dbo].[PR_Users_Update_GitHubUser]
(
        @UserPrincipalName varchar(100),
        @GithubUser varchar(100)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
UPDATE 
        [dbo].[Users]
   SET 
        [GithubUser] = @GithubUser,
        [Modified] = GETDATE(),
        [ModifiedBy] = @UserPrincipalName
 WHERE  
        [UserPrincipalName] = @UserPrincipalName
END
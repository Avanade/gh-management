CREATE PROCEDURE [dbo].[PR_Users_Update]
(
        @UserPrincipalName VARCHAR(100),
        @GivenName VARCHAR(50),
        @SurName VARCHAR(50),
        @JobTitle VARCHAR(50),
        @GitHubId VARCHAR(100),
        @GitHubUser VARCHAR(100),
        @ModifiedBy VARCHAR(100)
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
        [UserPrincipalName] = @UserPrincipalName,
        [GivenName] = @GivenName,
        [SurName] = @SurName,
        [JobTitle] = @JobTitle,
        [GitHubId] = @GitHubId,
        [GitHubUser] = @GithubUser,
        [Modified] = GETDATE(),
        [ModifiedBy] = @ModifiedBy
 WHERE  
        [UserPrincipalName] = @UserPrincipalName
END
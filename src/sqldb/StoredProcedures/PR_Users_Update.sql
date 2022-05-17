Create PROCEDURE [dbo].[PR_Users_Update]
(
        @UserPrincipalName varchar(100),
        @GivenName varchar(50),
        @SurName varchar(50),
        @JobTitle varchar(50),
        @GithubUser varchar(100),
        @ModifiedBy varchar(100)
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
        [GithubUser] = @GithubUser,
        [Modified] = GETDATE(),
        [ModifiedBy] = @ModifiedBy
 WHERE  
        [UserPrincipalName] = @UserPrincipalName
END
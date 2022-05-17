CREATE PROCEDURE PR_Users_Insert
(
			@UserPrincipalName varchar(100)
           ,@GivenName varchar(100)
           ,@SurName varchar(100)
           ,@JobTitle varchar(100)
           ,@GithubUser varchar(100)

)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here

INSERT INTO [dbo].[Users]
           (
           [UserPrincipalName],
           [GivenName],
           [SurName],
           [JobTitle],
           [GithubUser],
           [Created],
           [CreatedBy],
           [Modified],
           [ModifiedBy]
           )
     VALUES
           (
           @UserPrincipalName,
           @GivenName,
           @SurName,
           @JobTitle,
           @GithubUser,
           GETDATE(),
           @UserPrincipalName,
           GETDATE(),
           @UserPrincipalName
           )
END
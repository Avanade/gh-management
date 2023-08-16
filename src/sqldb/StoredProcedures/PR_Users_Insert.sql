CREATE PROCEDURE [dbo].[PR_Users_Insert]
(
			@UserPrincipalName VARCHAR(100)
           ,@Name VARCHAR(100)
           ,@GivenName VARCHAR(100) = NULL
           ,@SurName VARCHAR(100) = NULL
           ,@JobTitle VARCHAR(100) = NULL
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
	IF NOT EXISTS (SELECT UserPrincipalName FROM Users WHERE UserPrincipalName = @UserPrincipalName)
	BEGIN
		INSERT INTO [dbo].[Users]
			(
			[UserPrincipalName],
            [Name],
			[GivenName],
			[SurName],
			[JobTitle],
			[Created],
			[CreatedBy],
			[Modified],
			[ModifiedBy],
			[LastGithubLogin]
			)
		VALUES
			(
			@UserPrincipalName,
            @Name,
			@GivenName,
			@SurName,
			@JobTitle,
			GETDATE(),
			@UserPrincipalName,
			GETDATE(),
			@UserPrincipalName,
			DATEADD(day, -1, GETDATE())
			)
	END
END
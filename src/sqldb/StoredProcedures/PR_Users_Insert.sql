CREATE PROCEDURE [dbo].[PR_Users_Insert]
(
			@UserPrincipalName varchar(100)
           ,@GivenName varchar(100)
           ,@SurName varchar(100)
           ,@JobTitle varchar(100) = NULL
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
	IF NOT EXISTS (SELECT UserPrincipalName From Users WHERE UserPrincipalName = @UserPrincipalName)
	BEGIN
		INSERT INTO [dbo].[Users]
			(
			[UserPrincipalName],
			[GivenName],
			[SurName],
			[JobTitle],
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
			GETDATE(),
			@UserPrincipalName,
			GETDATE(),
			@UserPrincipalName
			)
	END
END
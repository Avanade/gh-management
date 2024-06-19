CREATE PROCEDURE [dbo].[usp_User_Insert]
  @UserPrincipalName [VARCHAR](100),
  @Name [VARCHAR](100),
  @GivenName [VARCHAR](100) = NULL,
  @SurName [VARCHAR](100) = NULL,
  @JobTitle [VARCHAR](100) = NULL
AS
BEGIN
	SET NOCOUNT ON;
  
	IF NOT EXISTS (SELECT [UserPrincipalName] FROM [dbo].[User] WHERE [UserPrincipalName] = @UserPrincipalName)
	BEGIN
		INSERT INTO [dbo].[User]
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
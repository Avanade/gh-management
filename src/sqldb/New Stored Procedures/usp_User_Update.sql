CREATE PROCEDURE [dbo].[usp_User_Update]
  @UserPrincipalName [VARCHAR](100),
  @GitHubId [VARCHAR](100),
  @GitHubUser [VARCHAR](100),
  @Force [BIT] = 0
AS
BEGIN
	SET NOCOUNT ON;

  DECLARE @LastGithubLogin [DATETIME]
	SET @LastGithubLogin = (SELECT [LastGithubLogin] FROM [dbo].[User] WHERE [UserPrincipalName] = @UserPrincipalName)

  IF EXISTS (
    SELECT [UserPrincipalName]
    FROM [dbo].[User]
    WHERE [UserPrincipalName] = @UserPrincipalName AND [GitHubId] IS NULL
  ) OR @Force = 1
  BEGIN
    UPDATE 
      [dbo].[User]
    SET
      [GitHubId] = @GitHubId,
      [GitHubUser] = @GitHubUser,
      [Modified] = GETDATE(),
      [ModifiedBy] = @UserPrincipalName,
      [LastGithubLogin] = GETDATE()
    WHERE [UserPrincipalName] = @UserPrincipalName

    SELECT CONVERT(BIT, 1) AS [IsValid], @GitHubId AS [GitHubId], @GitHubUser AS [GitHubUser], @LastGithubLogin AS [LastGithubLogin]
    RETURN 1
  END
  ELSE IF EXISTS(
    SELECT [UserPrincipalName]
    FROM [dbo].[User]
    WHERE [UserPrincipalName] = @UserPrincipalName AND [GitHubId] = @GitHubId AND [GitHubUser] != @GitHubUser
  )
  BEGIN
    UPDATE 
      [dbo].[User]
    SET
      [GitHubUser] = @GitHubUser,
      [Modified] = GETDATE(),
      [ModifiedBy] = @UserPrincipalName,
      [LastGithubLogin] = GETDATE()
    WHERE [UserPrincipalName] = @UserPrincipalName

    SELECT CONVERT(BIT, 1) AS [IsValid], @GitHubId AS [GitHubId], @GitHubUser AS [GitHubUser], @LastGithubLogin AS [LastGithubLogin]
    RETURN 1
  END
  ELSE
  BEGIN
    IF EXISTS (
      SELECT [UserPrincipalName]
      FROM [dbo].[User] 
      WHERE [UserPrincipalName] = @UserPrincipalName AND [GitHubId] = @GitHubId
    )
    BEGIN
      UPDATE [dbo].[User]
      SET [LastGithubLogin] = GETDATE()
      WHERE [UserPrincipalName] = @UserPrincipalName

      SELECT CONVERT(BIT, 1) AS [IsValid], @GitHubId AS [GitHubId], @GitHubUser AS [GitHubUser], @LastGithubLogin AS [LastGithubLogin]
      RETURN 1
    END
    ELSE
    BEGIN
      SELECT CONVERT(BIT, 0) AS [IsValid], @GitHubId AS [GitHubId], @GitHubUser AS [GitHubUser], @LastGithubLogin AS [LastGithubLogin]
      FROM [dbo].[User] WHERE [UserPrincipalName] = @UserPrincipalName
      RETURN 0
    END
  END
END
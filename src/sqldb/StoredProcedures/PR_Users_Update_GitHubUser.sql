CREATE PROCEDURE [dbo].[PR_Users_Update_GitHubUser]
(
        @UserPrincipalName VARCHAR(100),
        @GitHubId VARCHAR(100),
        @GitHubUser VARCHAR(100),
        @Force BIT = 0
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    DECLARE @LastGithubLogin DATETIME
	SET @LastGithubLogin = (SELECT LastGithubLogin FROM [dbo].[Users] WHERE UserPrincipalName=@UserPrincipalName)

    IF EXISTS (
        SELECT UserPrincipalName
        FROM Users
        WHERE
        UserPrincipalName = @UserPrincipalName
        AND GitHubId IS NULL
    ) OR @Force = 1
        BEGIN
            UPDATE 
                    [dbo].[Users]
            SET
                    [GitHubId] = @GitHubId,
                    [GitHubUser] = @GitHubUser,
                    [Modified] = GETDATE(),
                    [ModifiedBy] = @UserPrincipalName,
                    [LastGithubLogin] = GETDATE()
            WHERE  
                    [UserPrincipalName] = @UserPrincipalName

            SELECT CONVERT(BIT, 1) [IsValid], @GitHubId [GitHubId], @GitHubUser [GitHubUser], @LastGithubLogin [LastGithubLogin]
            RETURN 1
        END
    ELSE IF EXISTS(
        SELECT UserPrincipalName
        FROM Users
        WHERE
        UserPrincipalName = @UserPrincipalName
        AND GitHubId = @GitHubId AND GitHubUser != GitHubUser
    )
        BEGIN
            UPDATE 
                    [dbo].[Users]
            SET
                    [GitHubUser] = @GitHubUser,
                    [Modified] = GETDATE(),
                    [ModifiedBy] = @UserPrincipalName,
                    [LastGithubLogin] = GETDATE()
            WHERE  
                    [UserPrincipalName] = @UserPrincipalName

            SELECT CONVERT(BIT, 1) [IsValid], @GitHubId [GitHubId], @GitHubUser [GitHubUser], @LastGithubLogin [LastGithubLogin]
            RETURN 1
        END
    ELSE
        BEGIN
            IF EXISTS (
                SELECT UserPrincipalName
                FROM Users WHERE
                UserPrincipalName = @UserPrincipalName
                AND GitHubId = @GitHubId
            )
            BEGIN
                UPDATE 
                    [dbo].[Users]
                SET
                    [LastGithubLogin] = GETDATE()
                WHERE  
                    [UserPrincipalName] = @UserPrincipalName

                SELECT CONVERT(BIT, 1) [IsValid], @GitHubId [GitHubId], @GitHubUser [GitHubUser], @LastGithubLogin [LastGithubLogin]
                RETURN 1
            END
            ELSE
            BEGIN
                SELECT CONVERT(BIT, 0) [IsValid], GitHubId, GitHubUser, @LastGithubLogin [LastGithubLogin]
                FROM Users WHERE UserPrincipalName = @UserPrincipalName
                RETURN 0
            END
        END
END
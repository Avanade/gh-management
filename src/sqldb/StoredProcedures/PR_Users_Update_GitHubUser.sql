CREATE PROCEDURE [dbo].[PR_Users_Update_GitHubUser]
(
        @UserPrincipalName varchar(100),
        @GitHubId varchar(100),
        @GitHubUser varchar(100),
        @Force bit = 0
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

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
                    [ModifiedBy] = @UserPrincipalName
            WHERE  
                    [UserPrincipalName] = @UserPrincipalName

            SELECT CONVERT(BIT, 1) [IsValid], @GitHubId [GitHubId], @GitHubUser [GitHubUser]
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
                SELECT CONVERT(BIT, 1) [IsValid], @GitHubId [GitHubId], @GitHubUser [GitHubUser]
                RETURN 1
            END
            ELSE
            BEGIN
                SELECT CONVERT(BIT, 0) [IsValid], GitHubId, GitHubUser
                FROM Users WHERE UserPrincipalName = @UserPrincipalName
                RETURN 0
            END
        END
END
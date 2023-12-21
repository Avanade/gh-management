CREATE PROCEDURE [dbo].[PR_RepoOwners_Insert] 
(
    @ProjectId INT,
	@UserPrincipalName VARCHAR(100)
)
AS
BEGIN
    IF NOT EXISTS (
        SELECT * FROM RepoOwners WHERE 
            ProjectId = @ProjectId AND 
            UserPrincipalName = @UserPrincipalName
    )
    BEGIN
        INSERT INTO [dbo].[RepoOwners]
            ([ProjectId]
            ,[UserPrincipalName])
        VALUES
            (@ProjectId
            ,@UserPrincipalName)
    END
END
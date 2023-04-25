CREATE PROCEDURE [dbo].[PR_RepoOwners_Insert] 
(
    @ProjectId INT,
	@UserPrincipalName VARCHAR(100)
)
AS
BEGIN
    SET NOCOUNT ON
	INSERT INTO [dbo].[RepoOwners]
           ([ProjectId]
           ,[UserPrincipalName])
     VALUES
           (@ProjectId
           ,@UserPrincipalName)
END
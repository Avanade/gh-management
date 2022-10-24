CREATE PROCEDURE PR_CommunityMembers_Insert
	
	@CommunityId INT,
	@UserPrincipalName VARCHAR(100)

AS
IF NOT EXISTS (SELECT Id FROM CommunityMembers WHERE CommunityId = @CommunityId AND UserPrincipalName = @UserPrincipalName)
BEGIN
	INSERT INTO CommunityMembers (CommunityId, UserPrincipalName)
	VALUES (@CommunityId, @UserPrincipalName)
END
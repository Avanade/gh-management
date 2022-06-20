CREATE PROCEDURE PR_CommunityMembers_Remove
	
	@CommunityId INT,
	@UserPrincipalName VARCHAR(100)

AS

DELETE FROM CommunityMembers Where CommunityId = @CommunityId AND UserPrincipalName = @UserPrincipalName
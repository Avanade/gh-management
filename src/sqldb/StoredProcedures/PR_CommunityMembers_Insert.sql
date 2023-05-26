CREATE PROCEDURE PR_CommunityMembers_Insert
	
	@CommunityId INT,
	@UserPrincipalName VARCHAR(100)

AS
BEGIN
	INSERT INTO CommunityMembers (CommunityId, UserPrincipalName)
	VALUES (@CommunityId, @UserPrincipalName)
END
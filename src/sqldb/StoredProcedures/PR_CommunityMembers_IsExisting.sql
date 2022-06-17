CREATE PROCEDURE PR_CommunityMembers_IsExisting
	@CommunityId INT,
	@UserPrincipalName VARCHAR(100)
AS
IF EXISTS (
	SELECT UserPrincipalName FROM CommunityMembers WHERE CommunityId = @CommunityId AND UserPrincipalName = @UserPrincipalName
)
	BEGIN
		SELECT 1 [IsExisting]
		RETURN 1
	END
ELSE
	BEGIN
		SELECT 0 [IsExisting]
		RETURN 0
	END
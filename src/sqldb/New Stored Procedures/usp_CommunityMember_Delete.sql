CREATE PROCEDURE [dbo].[usp_CommunityMember_Delete]
	@CommunityId [INT],
	@UserPrincipalName [VARCHAR](100)
AS
BEGIN
  DELETE FROM [dbo].[CommunityMember] WHERE [CommunityId] = @CommunityId AND [UserPrincipalName] = @UserPrincipalName
END
CREATE PROCEDURE [dbo].[usp_CommunityMember_Insert]
  @CommunityId [INT],
  @UserPrincipalName [VARCHAR](100)
AS
BEGIN
  INSERT INTO [dbo].[CommunityMember]
    ([CommunityId], [UserPrincipalName])
  VALUES
    (@CommunityId, @UserPrincipalName)
END
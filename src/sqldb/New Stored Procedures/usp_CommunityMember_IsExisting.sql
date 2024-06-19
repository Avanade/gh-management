CREATE PROCEDURE [dbo].[usp_CommunityMember_IsExisting]
  @CommunityId [INT],
  @UserPrincipalName [VARCHAR](100)
AS
BEGIN
  IF EXISTS (
    SELECT [UserPrincipalName]
    FROM [dbo].[CommunityMember]
    WHERE [CommunityId] = @CommunityId AND [UserPrincipalName] = @UserPrincipalName
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
END
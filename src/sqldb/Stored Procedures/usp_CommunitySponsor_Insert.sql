CREATE PROCEDURE [dbo].[usp_CommunitySponsor_Insert]
		@CommunityId [INT],
		@UserPrincipalName [VARCHAR](100),
		@CreatedBy [VARCHAR](50)
AS
BEGIN
  SET NOCOUNT ON
	IF NOT EXISTS (SELECT * from [dbo].[CommunitySponsor] WHERE [UserPrincipalName] = @UserPrincipalName AND [CommunityId] = @CommunityId)
  BEGIN
    INSERT INTO [dbo].[CommunitySponsor]
    (
      [CommunityId],
      [UserPrincipalName],
      [Created],
      [CreatedBy],
      [Modified],
      [ModifiedBy]
    )
    VALUES 
    (
      @CommunityId,
      @UserPrincipalName,
      GETDATE(),
      @CreatedBy,
      GETDATE(),
      @CreatedBy
    )
  END
END
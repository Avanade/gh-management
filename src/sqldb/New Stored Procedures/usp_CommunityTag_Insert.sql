CREATE PROCEDURE [dbo].[usp_CommunityTag_Insert]
  @CommunityId [INT],
  @Tag [VARCHAR](20)
AS
BEGIN
  SET NOCOUNT ON

  INSERT INTO [dbo].[CommunityTag]
  (
    [CommunityId],
    [Tag]
  )
  VALUES
  (
      @CommunityId,
      @Tag
  )
END
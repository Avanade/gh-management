CREATE PROCEDURE [dbo].[usp_CommunityTag_Select_ByCommunityId]
  @CommunityId [INT]
AS
BEGIN
  SET NOCOUNT ON

  SELECT
    [Id],
    [CommunityId],
    [Tag]
  FROM [dbo].[CommunityTag]
  WHERE [CommunityId] = @CommunityId
END
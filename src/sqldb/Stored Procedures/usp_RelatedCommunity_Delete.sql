CREATE PROCEDURE [dbo].[usp_RelatedCommunity_Delete]
  @ParentCommunityId [INT]
AS
BEGIN
  DELETE  [dbo].[RelatedCommunity]
  WHERE [ParentCommunityId] = @ParentCommunityId
END

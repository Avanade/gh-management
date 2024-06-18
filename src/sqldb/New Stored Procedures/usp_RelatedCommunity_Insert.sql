CREATE PROCEDURE [dbo].[usp_RelatedCommunity_Insert]
  @ParentCommunityId [INT],
  @RelatedCommunityId [INT]
AS
BEGIN
	IF (@ParentCommunityId != @RelatedCommunityId)
	BEGIN
		SET NOCOUNT ON
		INSERT INTO [dbo].[RelatedCommunity]
    (
      [ParentCommunityId],
      [RelatedCommunityId]
    )
		 VALUES
    (
      @ParentCommunityId,
			@RelatedCommunityId
    )
	END
END

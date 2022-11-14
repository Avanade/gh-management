create PROCEDURE [dbo].[PR_RelatedCommunities_Delete]
(
 @ParentCommunityId int
 
 )
AS
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.

delete  [dbo].[RelatedCommunities]
  where 
		ParentCommunityId = @ParentCommunityId
	 

END
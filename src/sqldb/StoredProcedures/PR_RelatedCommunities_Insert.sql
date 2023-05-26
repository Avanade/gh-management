create PROCEDURE [dbo].[PR_RelatedCommunities_Insert]
(
    -- Add the parameters for the stored procedure here
   @ParentCommunityId int,
   @RelatedCommunityId int
)
AS
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.
	if (@ParentCommunityId != @RelatedCommunityId)
	begin
		SET NOCOUNT ON
		INSERT INTO [dbo].[RelatedCommunities]
			   ([ParentCommunityId]
			   ,[RelatedCommunityId])
		 VALUES
			   (@ParentCommunityId
			   ,@RelatedCommunityId)
	end
END
 
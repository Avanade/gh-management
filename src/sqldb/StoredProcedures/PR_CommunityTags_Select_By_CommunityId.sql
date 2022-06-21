
create PROCEDURE  [dbo].[PR_CommunityTags_Select_By_CommunityId]
 @CommunityId int
AS
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.
    SET NOCOUNT ON

    -- Insert statements for procedure here
 SELECT [Id]
      ,[CommunityId]
      ,[Tag]
  FROM [dbo].[CommunityTags]
  where [CommunityId] = @CommunityId
END

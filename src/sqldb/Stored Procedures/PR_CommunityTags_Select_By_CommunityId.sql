
CREATE PROCEDURE  [dbo].[PR_CommunityTags_Select_By_CommunityId]
 @CommunityId INT
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
  WHERE [CommunityId] = @CommunityId
END

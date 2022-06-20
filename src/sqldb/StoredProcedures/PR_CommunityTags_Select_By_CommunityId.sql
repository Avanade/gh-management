/****** Object:  StoredProcedure [dbo].[PR_CommunitySponsors_Select_By_CommunityId]    Script Date: 6/20/2022 2:39:44 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:      <Author, , Name>
-- Create Date: <Create Date, , >
-- Description: <Description, , >
-- =============================================
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

/****** Object:  StoredProcedure [dbo].[PR_CommunityActivities_Select_ByOffsetAndFilter]    Script Date: 6/24/2022 11:37:16 AM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[PR_CommunityActivities_TotalCount]
AS
BEGIN
    SET NOCOUNT ON
	SELECT COUNT(Id) AS Total FROM CommunityActivities
END
/****** Object:  StoredProcedure [dbo].[PR_CommunityActivities_TotalCount]    Script Date: 29/06/2022 19:53:39 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[PR_CommunityActivities_TotalCount_ByCreatedBy] (
	@CreatedBy varchar(50)
)
AS
BEGIN
    SET NOCOUNT ON
	SELECT COUNT(Id) AS Total FROM CommunityActivities WHERE CreatedBy = @CreatedBy
END
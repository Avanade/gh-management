CREATE PROCEDURE [dbo].[PR_CommunityActivities_TotalCount]
AS
BEGIN
    SET NOCOUNT ON
	SELECT COUNT(Id) AS Total FROM CommunityActivities
END
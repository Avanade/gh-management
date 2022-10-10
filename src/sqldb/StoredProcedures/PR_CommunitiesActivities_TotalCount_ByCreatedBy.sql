CREATE PROCEDURE [dbo].[PR_CommunityActivities_TotalCount_ByCreatedBy] (
	@CreatedBy VARCHAR(50)
)
AS
BEGIN
    SET NOCOUNT ON
	SELECT COUNT(Id) AS Total FROM CommunityActivities WHERE CreatedBy = @CreatedBy
END
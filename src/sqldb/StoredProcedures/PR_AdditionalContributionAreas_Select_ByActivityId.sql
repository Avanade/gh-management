CREATE PROCEDURE [dbo].[PR_AdditionalContributionAreas_Select_ByActivityId]
(
	@ActivityId INT
)
AS
BEGIN
    SET NOCOUNT ON
    SELECT 
	* 
FROM 
		CommunityActivitiesContributionAreas AS caca 
	JOIN
		ContributionAreas AS ca ON caca.ContributionAreaId = ca.Id
	WHERE
		caca.CommunityActivityId = 43 AND caca.IsPrimary = 0
END
CREATE PROCEDURE [dbo].[PR_ContributionAreas_TotalCount]
AS
BEGIN
	SELECT COUNT(Id) AS 'Total' FROM [dbo].[ContributionAreas]
END
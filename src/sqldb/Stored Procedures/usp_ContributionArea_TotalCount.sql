CREATE PROCEDURE [dbo].[usp_ContributionArea_TotalCount]
AS
BEGIN
	SELECT COUNT(*) AS [Total] FROM [dbo].[ContributionArea]
END
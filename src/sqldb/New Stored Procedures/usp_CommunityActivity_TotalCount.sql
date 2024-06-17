CREATE PROCEDURE [dbo].[usp_CommunityActivity_TotalCount]
AS
BEGIN
  SET NOCOUNT ON
	SELECT COUNT(*) AS [Total] FROM [dbo].[CommunityActivity]
END
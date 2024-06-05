CREATE PROCEDURE [dbo].[usp_CommunityActivity_TotalCount_ByOption]
AS
BEGIN
  SET NOCOUNT ON
	SELECT COUNT(*) AS [Total] FROM [dbo].[CommunityActivity]
END
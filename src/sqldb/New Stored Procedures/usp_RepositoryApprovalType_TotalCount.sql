CREATE PROCEDURE [dbo].[usp_RepositoryApprovalType_TotalCount]
AS
BEGIN
	SELECT COUNT(*) AS [Total] FROM [dbo].[RepositoryApprovalType] WHERE [IsArchived] = 0
END
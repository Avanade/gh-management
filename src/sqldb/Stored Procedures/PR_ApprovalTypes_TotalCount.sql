CREATE PROCEDURE [dbo].[PR_ApprovalTypes_TotalCount]
AS
BEGIN
	SELECT COUNT(Id) AS 'Total' FROM [dbo].[ApprovalTypes] WHERE IsArchived = 0
END
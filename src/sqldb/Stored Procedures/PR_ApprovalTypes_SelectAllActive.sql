CREATE PROCEDURE [dbo].[PR_ApprovalTypes_SelectAllActive]
AS
BEGIN
	SELECT * FROM [dbo].[ApprovalTypes]
	WHERE [IsActive] = 1 AND [IsArchived] = 0
END
CREATE PROCEDURE [dbo].[PR_ApprovalTypes_Select_ById] 
(
	@Id INT
)
AS
BEGIN
	SELECT * FROM [dbo].[ApprovalTypes] WHERE Id = @Id
END
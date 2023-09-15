CREATE PROCEDURE [dbo].[PR_Approvers_Select_ByApprovalTypeId] 
(
	@Id INT
)
AS
BEGIN
	SELECT * FROM [dbo].[Approvers] WHERE Id = @Id
END
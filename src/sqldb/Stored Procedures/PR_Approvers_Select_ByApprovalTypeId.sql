CREATE PROCEDURE [dbo].[PR_Approvers_Select_ByApprovalTypeId] 
(
	@ApprovalTypeId INT
)
AS
BEGIN
	SELECT * FROM [dbo].[Approvers] WHERE ApprovalTypeId = @ApprovalTypeId
END
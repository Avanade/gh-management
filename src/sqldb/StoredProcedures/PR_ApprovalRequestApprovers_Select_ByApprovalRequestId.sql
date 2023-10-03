CREATE PROCEDURE [dbo].[PR_ApprovalRequestApprovers_Select_ByApprovalRequestId] 
(
	@ApprovalRequestId INT
)
AS
BEGIN
	SELECT * FROM [dbo].[ApprovalRequestApprovers] WHERE ApprovalRequestId = @ApprovalRequestId
END
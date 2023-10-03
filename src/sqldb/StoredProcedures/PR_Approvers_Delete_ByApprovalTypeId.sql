CREATE PROCEDURE  [dbo].[PR_Approvers_Delete_ByApprovalTypeId]
(
    @ApprovalTypeId INT
)
AS
BEGIN   
    DELETE FROM Approvers WHERE ApprovalTypeId = @ApprovalTypeId
END

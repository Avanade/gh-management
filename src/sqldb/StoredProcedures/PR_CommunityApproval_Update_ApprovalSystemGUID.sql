Create PROCEDURE [dbo].[PR_CommunityApproval_Update_ApprovalSystemGUID]
    @Id INT,
    @ApprovalSystemGUID UNIQUEIDENTIFIER

AS

UPDATE CommunityApprovals
SET
    ApprovalStatusId = 2,
    ApprovalSystemGUID = @ApprovalSystemGUID,
    ApprovalSystemDateSent = GETDATE(),
    Modified = GETDATE()
WHERE Id = @Id
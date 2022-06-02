CREATE PROCEDURE PR_ProjectsApproval_Update_ApprovalSystemGUID
    @Id INT,
    @ApprovalSystemGUID UNIQUEIDENTIFIER

AS

UPDATE ProjectApprovals
SET
    ApprovalStatusId = 2,
    ApprovalSystemGUID = @ApprovalSystemGUID,
    ApprovalSystemDateSent = GETDATE(),
    Modified = GETDATE()
WHERE Id = @Id
CREATE PROCEDURE [dbo].[PR_Organizations_UpdateStatus]
    @Id INT,
    @ApprovalStatusId INT
AS
UPDATE CommunityApprovals
SET
    ApprovalStatusId = @ApprovalStatusId
WHERE Id = @Id
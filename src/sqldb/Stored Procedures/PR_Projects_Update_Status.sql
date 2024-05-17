CREATE PROCEDURE [dbo].[PR_Projects_Update_Status]
  @ProjectId INT

AS

IF EXISTS (SELECT ApprovalStatusId FROM ProjectApprovals WHERE ProjectId = @ProjectId AND ApprovalStatusId = 3) -- IF REJECTED
BEGIN
	UPDATE Projects SET ApprovalStatusId = 3 WHERE Id = @ProjectId
END

ELSE IF NOT EXISTS (SELECT ApprovalStatusId FROM ProjectApprovals WHERE ProjectId = @ProjectId AND ApprovalStatusId NOT IN (3,5)) -- EVERYONE HAS RESPONDED
BEGIN
	UPDATE Projects SET ApprovalStatusId = 5 WHERE Id = @ProjectId
END
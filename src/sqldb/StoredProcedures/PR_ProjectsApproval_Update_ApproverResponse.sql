CREATE PROCEDURE [dbo].[PR_ProjectsApproval_Update_ApproverResponse]
(
  @ApprovalSystemGUID UNIQUEIDENTIFIER,
  @ApprovalStatusId INT,
  @ApprovalRemarks VARCHAR(255),
  @ApprovalDate DATETIME
)
AS
BEGIN

UPDATE
	[dbo].[ProjectApprovals]
  SET
    [ApprovalStatusId] = @ApprovalStatusId,
    [ApprovalRemarks] = @ApprovalRemarks,
    [ModifiedBy] = [ApproverUserPrincipalName],
    [Modified] = GETDATE(),
    [ApprovalDate] = convert(DATETIME, @ApprovalDate)
  WHERE
    [ApprovalSystemGUID] = @ApprovalSystemGUID
END

DECLARE @ProjectId INT
SELECT @ProjectId = ProjectId FROM ProjectApprovals WHERE [ApprovalSystemGUID] = @ApprovalSystemGUID

EXEC PR_Projects_Update_Status @ProjectId
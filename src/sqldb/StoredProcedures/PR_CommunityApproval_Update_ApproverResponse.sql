CREATE PROCEDURE [dbo].[PR_CommunityApproval_Update_ApproverResponse]
(
  @ApprovalSystemGUID UNIQUEIDENTIFIER,
  @ApprovalStatusId INT,
  @ApprovalRemarks VARCHAR(255),
  @ApprovalDate DATETIME
)
AS
BEGIN
  -- SET NOCOUNT ON added to prevent extra result sets from
  -- interfering with SELECT statements.
  SET NOCOUNT ON

UPDATE
	[dbo].[CommunityApprovals]
  SET
    [ApprovalStatusId] = @ApprovalStatusId,
    [ApprovalRemarks] = @ApprovalRemarks,
    [ModifiedBy] = [ApproverUserPrincipalName],
    [Modified] = GETDATE(),
    [ApprovalDate] = convert(DATETIME, @ApprovalDate)
  WHERE
    [ApprovalSystemGUID] = @ApprovalSystemGUID
END

DECLARE @RequestId INT
SELECT @RequestId = Id FROM CommunityApprovals WHERE [ApprovalSystemGUID] = @ApprovalSystemGUID

DECLARE @CommunityId INT
SELECT @CommunityId = CommunityId FROM CommunityApprovalRequests WHERE [RequestId] = @RequestId

EXEC PR_Communities_Update_Status @CommunityId
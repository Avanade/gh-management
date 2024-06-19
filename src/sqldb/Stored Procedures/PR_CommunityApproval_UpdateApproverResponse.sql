CREATE PROCEDURE [dbo].[PR_CommunityApproval_UpdateApproverResponse]
(
  @ApprovalSystemGUID UNIQUEIDENTIFIER,
  @ApprovalStatusId INT,
  @ApprovalRemarks VARCHAR(255),
  @ApprovalDate DATETIME,
  @Approver VARCHAR(100)
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
    [ApprovalSystemGUID] = @ApprovalSystemGUID AND
    [ApproverUserPrincipalName] = @Approver;

  UPDATE
  	[dbo].[CommunityApprovals]
  SET
    [ApprovalStatusId] = 7,
    [Modified] = GETDATE()
  WHERE
    [ApprovalSystemGUID] = @ApprovalSystemGUID AND
    [ApproverUserPrincipalName] != @Approver;
END
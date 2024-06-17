CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Update_ApproverResponse]
  @ApprovalSystemGUID [UNIQUEIDENTIFIER],
  @ApprovalStatusId [INT],
  @ApprovalRemarks [VARCHAR](255),
  @ApprovalDate [DATETIME],
  @Approver [VARCHAR](100)
AS
BEGIN
  SET NOCOUNT ON

  UPDATE
    [dbo].[ApprovalRequest]
  SET
    [ApprovalStatusId] = @ApprovalStatusId,
    [ApprovalRemarks] = @ApprovalRemarks,
    [ModifiedBy] = [ApproverUserPrincipalName],
    [Modified] = GETDATE(),
    [ApprovalDate] = convert([DATETIME], @ApprovalDate)
  WHERE
    [ApprovalSystemGUID] = @ApprovalSystemGUID AND
    [ApproverUserPrincipalName] = @Approver;

  UPDATE
  	[dbo].[ApprovalRequest]
  SET
    [ApprovalStatusId] = 7,
    [Modified] = GETDATE()
  WHERE
    [ApprovalSystemGUID] = @ApprovalSystemGUID AND
    [ApproverUserPrincipalName] != @Approver;
END
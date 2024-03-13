CREATE PROCEDURE [dbo].[PR_CommunityApproval_UpdateOrganization_ApproverResponse]
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
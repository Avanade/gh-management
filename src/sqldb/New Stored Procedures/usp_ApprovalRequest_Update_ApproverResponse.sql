CREATE PROCEDURE [dbo].[usp_ApprovalRequest_Update_ApproverResponse]
  @ApprovalSystemGUID [UNIQUEIDENTIFIER],
  @ApprovalStatusId [INT],
  @ApprovalRemarks [VARCHAR](255),
  @ApprovalDate [DATETIME]
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
    [ApprovalDate] = convert(DATETIME, @ApprovalDate)
  WHERE
    [ApprovalSystemGUID] = @ApprovalSystemGUID;
    
  DECLARE @CommunityId INT

  SELECT 
    @CommunityId = [CAR].[CommunityId] 
  FROM
    [dbo].[CommunityApprovalRequest] AS [CAR]
    INNER JOIN [dbo].[ApprovalRequest] AS [AR] ON [AR].[Id] = [CAR].[ApprovalRequestId]
  WHERE
    [AR].[ApprovalSystemGUID] = @ApprovalSystemGUID

  -- EXEC usp_Community_Update_Status @CommunityId
END
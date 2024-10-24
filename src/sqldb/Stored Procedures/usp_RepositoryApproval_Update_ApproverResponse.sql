CREATE PROCEDURE [dbo].[usp_RepositoryApproval_Update_ApproverResponse]
  @ApprovalSystemGUID [UNIQUEIDENTIFIER],
  @ApprovalStatusId [INT],
  @ApprovalRemarks [VARCHAR](255),
  @ApprovalDate [DATETIME],
  @RespondedBy [VARCHAR](100)
AS
BEGIN
  UPDATE [dbo].[RepositoryApproval]
  SET
    [ApprovalStatusId] = @ApprovalStatusId,
    [ApprovalRemarks] = @ApprovalRemarks,
    [ModifiedBy] = @RespondedBy,
    [Modified] = GETDATE(),
    [ApprovalDate] = CONVERT(DATETIME, @ApprovalDate),
    [RespondedBy] = @RespondedBy
  WHERE
    [ApprovalSystemGUID] = @ApprovalSystemGUID

  DECLARE @RepositoryId [INT]
  SELECT @RepositoryId = [RepositoryId] FROM [dbo].[RepositoryApproval] WHERE [ApprovalSystemGUID] = @ApprovalSystemGUID

  EXEC [dbo].[usp_Repository_Update_Status] @RepositoryId
END
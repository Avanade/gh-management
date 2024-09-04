CREATE PROCEDURE [dbo].[usp_RepositoryApproval_Update_ApprovalSystemGUID]
  @Id [INT],
  @ApprovalSystemGUID [UNIQUEIDENTIFIER]
AS
BEGIN
  UPDATE [dbo].[RepositoryApproval]
  SET
    [ApprovalStatusId] = 2,
    [ApprovalSystemGUID] = @ApprovalSystemGUID,
    [ApprovalSystemDateSent] = GETDATE(),
    [Modified] = GETDATE()
  WHERE [Id] = @Id
END

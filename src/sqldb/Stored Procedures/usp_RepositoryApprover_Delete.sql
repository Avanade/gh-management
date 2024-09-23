CREATE PROCEDURE [dbo].[usp_RepositoryApprover_Delete]
  @RepositoryApprovalTypeId [INT]
AS
BEGIN
  SET NOCOUNT ON

  DELETE FROM [dbo].[RepositoryApprover] WHERE [RepositoryApprovalTypeId] = @RepositoryApprovalTypeId
END
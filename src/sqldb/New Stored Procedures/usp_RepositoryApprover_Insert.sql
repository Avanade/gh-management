CREATE PROCEDURE [dbo].[usp_RepositoryApprover_Insert]
  @RepositoryApprovalTypeId [INT],
  @ApproverUserPrincipalName [VARCHAR](100)
AS
BEGIN
  SET NOCOUNT ON

  IF NOT EXISTS (
    SELECT * FROM [dbo].[RepositoryApprover]
    WHERE [RepositoryApprovalTypeId] = @RepositoryApprovalTypeId AND [ApproverUserPrincipalName] = @ApproverUserPrincipalName
  )
    BEGIN
    INSERT INTO [dbo].[RepositoryApprover]
    (
      [RepositoryApprovalTypeId],
      [ApproverUserPrincipalName]
    )
    VALUES
    (
        @RepositoryApprovalTypeId,
        @ApproverUserPrincipalName
    )
  END
END
CREATE PROCEDURE [dbo].[usp_ApprovalRequestApprover_Select_ByApprovalRequestId]
	@RepositoryApprovalId [INT]
AS
BEGIN
	SELECT
    [RepositoryApprovalId],
    [ApproverUserPrincipalName]
  FROM 
    [dbo].[ApprovalRequestApprover] 
  WHERE [RepositoryApprovalId] = @RepositoryApprovalId
END
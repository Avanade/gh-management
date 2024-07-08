CREATE PROCEDURE [dbo].[usp_RepositoryApprover_Select_ByApprovalTypeId]
	@RepositoryApprovalTypeId [INT]
AS
BEGIN
	SELECT
		[R].[RepositoryApprovalTypeId],
		[R].[ApproverUserPrincipalName],
		[U].[Name] AS [ApproverName]
	FROM [dbo].[RepositoryApprover] AS [R]
	INNER JOIN [dbo].[User] AS [U] ON [R].[ApproverUserPrincipalName] = [U].[UserPrincipalName]
	WHERE [RepositoryApprovalTypeId] = @RepositoryApprovalTypeId
END
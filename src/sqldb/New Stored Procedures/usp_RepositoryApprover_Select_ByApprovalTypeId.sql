CREATE PROCEDURE [dbo].[usp_RepositoryApprover_Select_ByApprovalTypeId]
	@RepositoryApprovalTypeId [INT]
AS
BEGIN
	SELECT * FROM [dbo].[RepositoryApprover] WHERE [RepositoryApprovalTypeId] = @RepositoryApprovalTypeId
END
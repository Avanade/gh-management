CREATE PROCEDURE [dbo].[usp_RepositoryApprovalType_Select_Active]
AS
BEGIN
	SELECT
    [Id],
    [Name],
    [ApproverUserPrincipalName],
    [IsArchived],
    [IsActive],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy]
  FROM [dbo].[RepositoryApprovalType]
	WHERE [IsActive] = 1 AND [IsArchived] = 0
END
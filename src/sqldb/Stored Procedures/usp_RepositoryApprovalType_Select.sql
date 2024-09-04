CREATE PROCEDURE [dbo].[usp_RepositoryApprovalType_Select]
AS
BEGIN
	SELECT
    [Id],
    [Name],
    [IsArchived],
    [IsActive],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy]
  FROM [dbo].[RepositoryApprovalType]
END
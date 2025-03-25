CREATE PROCEDURE [dbo].[usp_RepositoryApprovalType_Select_ById]
	@Id [INT]
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
	FROM [dbo].[RepositoryApprovalType] WHERE [Id] = @Id
END
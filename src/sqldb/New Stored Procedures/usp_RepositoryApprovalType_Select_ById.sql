CREATE PROCEDURE [dbo].[usp_RepositoryApprovalType_Select_ById]
	@Id [INT]
AS
BEGIN
	SELECT * FROM [dbo].[RepositoryApprovalType] WHERE [Id] = @Id
END
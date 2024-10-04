CREATE PROCEDURE [dbo].[usp_Repository_Delete]
	@RepositoryId [INT]
AS
BEGIN
    DELETE FROM [dbo].[RepositoryTopic] WHERE [RepositoryId] = @RepositoryId;
    DELETE FROM [dbo].[RepositoryOwner] WHERE [RepositoryId] = @RepositoryId;
    DELETE FROM [dbo].[ApprovalRequestApprover] WHERE [RepositoryApprovalId] IN (SELECT [Id] FROM [dbo].[RepositoryApproval] WHERE [RepositoryId] = @RepositoryId);
    DELETE FROM [dbo].[RepositoryApproval] WHERE [RepositoryId] = @RepositoryId;
    DELETE FROM [dbo].[Repository] WHERE [Id] = @RepositoryId;
END
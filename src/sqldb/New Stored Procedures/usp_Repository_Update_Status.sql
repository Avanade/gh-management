CREATE PROCEDURE [dbo].[usp_Repository_Update_Status]
  @RepositoryId [INT]
AS
BEGIN
  -- IF REJECTED
  IF EXISTS (SELECT [ApprovalStatusId] FROM [dbo].[RepositoryApproval] WHERE [RepositoryId] = @RepositoryId AND [ApprovalStatusId] = 3)
  BEGIN
    UPDATE [dbo].[Repository] SET [ApprovalStatusId] = 3 WHERE [Id] = @RepositoryId
  END
  -- EVERYONE HAS RESPONDED
  ELSE IF NOT EXISTS (SELECT [ApprovalStatusId] FROM [dbo].[RepositoryApproval] WHERE [RepositoryId] = @RepositoryId AND [ApprovalStatusId] NOT IN (3,5))
  BEGIN
    UPDATE [dbo].[Repository] SET [ApprovalStatusId] = 5 WHERE [Id] = @RepositoryId
  END
END
ALTER PROCEDURE [dbo].[PR_Projects_Delete_ById]
(
	@Id INT
)
AS
BEGIN
    DELETE FROM [dbo].[RepoTopics] WHERE ProjectId = @Id;
    DELETE FROM [dbo].[RepoOwners] WHERE ProjectId = @Id;
    DELETE FROM [dbo].[ApprovalRequestApprovers] WHERE ApprovalRequestId IN (SELECT Id FROM ProjectApprovals WHERE ProjectId = @Id)
    DELETE FROM [dbo].[ProjectApprovals] WHERE ProjectId = @Id;
    DELETE FROM [dbo].[Projects] WHERE Id = @Id;

END
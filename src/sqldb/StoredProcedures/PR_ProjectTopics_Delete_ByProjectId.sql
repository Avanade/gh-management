CREATE PROCEDURE [dbo].[PR_RepoTopics_Delete_ByProjectId] 
(
	@ProjectId INT
)
AS
BEGIN
	DELETE FROM [dbo].[RepoTopics] WHERE ProjectId = @ProjectId
END
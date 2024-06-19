CREATE PROCEDURE [dbo].[PR_RepoTopics_Insert] 
(
	@Topic VARCHAR(100),
	@ProjectId INT
)
AS
BEGIN
	INSERT INTO [dbo].[RepoTopics]
           (Topic, ProjectId)
     VALUES
           (@Topic, @ProjectId)
END
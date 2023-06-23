CREATE TABLE [dbo].[RepoTopics]
(
	[Topic] VARCHAR(100) NOT NULL,
    [ProjectId] INT NOT NULL,
    CONSTRAINT [FK_RepoTags_Project] FOREIGN KEY (ProjectId) REFERENCES Projects(Id),
	PRIMARY KEY (Topic, ProjectId)
)
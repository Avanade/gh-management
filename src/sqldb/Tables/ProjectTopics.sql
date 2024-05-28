CREATE TABLE [dbo].[RepoTopics] (
    [Topic] [VARCHAR](100) NOT NULL,
    [ProjectId] [INT] NOT NULL,
    CONSTRAINT [PK_RepoTopics] PRIMARY KEY ([Topic], [ProjectId]),
    CONSTRAINT [FK_RepoTags_Project] FOREIGN KEY ([ProjectId]) REFERENCES [dbo].[Projects]([Id])
)

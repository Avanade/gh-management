CREATE TABLE [dbo].[RepositoryTopic] (
    [Topic] [VARCHAR](100) NOT NULL,
    [ProjectId] [INT] NOT NULL,
    CONSTRAINT [PK_RepositoryTopic] PRIMARY KEY ([Topic], [ProjectId]),
    CONSTRAINT [FK_RepositoryTopic_Repository] FOREIGN KEY ([ProjectId]) REFERENCES [dbo].[Repository]([Id])
)

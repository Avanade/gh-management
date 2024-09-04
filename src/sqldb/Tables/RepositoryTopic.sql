CREATE TABLE [dbo].[RepositoryTopic] (
    [Topic] [VARCHAR](100) NOT NULL,
    [RepositoryId] [INT] NOT NULL,
    CONSTRAINT [PK_RepositoryTopic] PRIMARY KEY ([Topic], [RepositoryId]),
    CONSTRAINT [FK_RepositoryTopic_Repository] FOREIGN KEY ([RepositoryId]) REFERENCES [dbo].[Repository]([Id])
)

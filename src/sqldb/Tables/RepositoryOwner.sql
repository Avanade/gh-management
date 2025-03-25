CREATE TABLE [dbo].[RepositoryOwner] (
    [RepositoryId] [INT] NOT NULL,
    [UserPrincipalName] [VARCHAR](100) NOT NULL,
    CONSTRAINT [PK_RepositoryOwner] PRIMARY KEY ([RepositoryId], [UserPrincipalName]),
    CONSTRAINT [FK_RepositoryOwner_Repository] FOREIGN KEY ([RepositoryId]) REFERENCES [dbo].[Repository]([Id])
)

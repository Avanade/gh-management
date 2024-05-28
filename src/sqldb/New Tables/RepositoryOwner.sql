CREATE TABLE [dbo].[RepositoryOwner] (
    [ProjectId] [INT] NOT NULL,
    [UserPrincipalName] [VARCHAR](100) NOT NULL,
    CONSTRAINT [PK_RepositoryOwner] PRIMARY KEY ([ProjectId], [UserPrincipalName]),
    CONSTRAINT [FK_RepositoryOwner_Repository] FOREIGN KEY ([ProjectId]) REFERENCES [dbo].[Repository]([Id])
)
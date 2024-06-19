CREATE TABLE [dbo].[RepoOwners] (
    [ProjectId] [INT] NOT NULL,
    [UserPrincipalName] [VARCHAR](100) NOT NULL,
    CONSTRAINT [PK_RepoOwner] PRIMARY KEY ([ProjectId], [UserPrincipalName]),
    CONSTRAINT [FK_RepoOwners_Projects] FOREIGN KEY ([ProjectId]) REFERENCES [dbo].[Projects]([Id])
)
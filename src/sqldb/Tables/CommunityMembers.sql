CREATE TABLE [dbo].[CommunityMembers] (
    [CommunityId] [INT] NOT NULL,
    [UserPrincipalName] [VARCHAR](100) NOT NULL,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
    [Modified] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [ModifiedBy] [VARCHAR](100) NULL,
    CONSTRAINT [PK_CommunityMembers] PRIMARY KEY ([CommunityId], [UserPrincipalName]),
    CONSTRAINT [FK_CommunityMembers_Communities] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Communities]([Id])
)

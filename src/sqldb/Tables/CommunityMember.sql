CREATE TABLE [dbo].[CommunityMember] (
    [CommunityId] [INT] NOT NULL,
    [UserPrincipalName] [VARCHAR](100) NOT NULL,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
    [Modified] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [ModifiedBy] [VARCHAR](100) NULL,
    CONSTRAINT [PK_CommunityMember] PRIMARY KEY ([CommunityId], [UserPrincipalName]),
    CONSTRAINT [FK_CommunityMember_Community] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Community]([Id])
)

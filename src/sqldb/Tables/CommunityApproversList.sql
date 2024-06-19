CREATE TABLE [dbo].[CommunityApproversList] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [ApproverUserPrincipalName] [VARCHAR](100) NOT NULL,
    [Category] [VARCHAR](100) NOT NULL DEFAULT 'community',
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
    [Modified] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [ModifiedBy] [VARCHAR](100) NULL,
    [Disabled] [BIT] NULL,
    CONSTRAINT [FK_CommunityApproversList_Users] FOREIGN KEY ([ApproverUserPrincipalName]) REFERENCES [dbo].[Users]([UserPrincipalName]),
    CONSTRAINT [AK_ApproverUserPrincipalName_Categoy] UNIQUE ([ApproverUserPrincipalName], [Category])
)

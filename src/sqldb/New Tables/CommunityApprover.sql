CREATE TABLE [dbo].[CommunityApprover] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [ApproverUserPrincipalName] [VARCHAR](100) NOT NULL,
    [GuidanceCategory] [VARCHAR](100) NOT NULL DEFAULT 'community',
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
    [Modified] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [ModifiedBy] [VARCHAR](100) NULL,
    [Disabled] [BIT] NULL,
    CONSTRAINT [FK_CommunityApprover_User] FOREIGN KEY ([ApproverUserPrincipalName]) REFERENCES [dbo].[User]([UserPrincipalName]),
    CONSTRAINT [AK_ApproverUserPrincipalName_GuidanceCategory] UNIQUE ([ApproverUserPrincipalName], [GuidanceCategory])
)

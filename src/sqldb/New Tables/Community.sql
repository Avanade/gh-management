CREATE TABLE [dbo].[Community] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [Name] [VARCHAR](50) NOT NULL,
    [Url] [VARCHAR](255) NULL,
    [Description] [VARCHAR](255) NULL,
    [Notes] [VARCHAR](255) NULL,
    [TradeAssocId] [VARCHAR](50) NULL,
    [IsExternal] [BIT] NOT NULL DEFAULT 0,
    [ApprovalStatusId] [INT] NOT NULL DEFAULT 1,
    [OnBoardingInstructions] [VARCHAR](MAX) NULL,
    [CommunityType] [VARCHAR](10) NULL,
    [ChannelId] [VARCHAR](100) NULL,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
    [Modified] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [ModifiedBy] [VARCHAR](100) NULL,
    CONSTRAINT [FK_Community_ApprovalStatus] FOREIGN KEY ([ApprovalStatusId]) REFERENCES [dbo].[ApprovalStatus]([Id])
)

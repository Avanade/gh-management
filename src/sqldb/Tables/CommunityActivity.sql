CREATE TABLE [dbo].[CommunityActivity] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [CommunityId] [INT] NOT NULL,
    [Date] [DATETIME],
    [Name] [VARCHAR](255) NOT NULL,
    [ActivityTypeId] [INT] NOT NULL,
    [Url] [VARCHAR](255) NULL,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
    [Modified] [DATETIME] NULL,
    [ModifiedBy] [VARCHAR](100) NULL,
    CONSTRAINT [FK_CommunityActivity_Community] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Community]([Id]),
    CONSTRAINT [FK_CommunityActivity_ActivityType] FOREIGN KEY ([ActivityTypeId]) REFERENCES [dbo].[ActivityType]([Id]),
)

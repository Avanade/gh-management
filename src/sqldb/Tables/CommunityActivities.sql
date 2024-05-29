CREATE TABLE [dbo].[CommunityActivities] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [CommunityId] [INT] NOT NULL,
    [Date] [DATETIME],
    [Name] [VARCHAR](255) NOT NULL,
    [ActivityTypeId] [INT] NOT NULL,
    [Url] [VARCHAR](255) NULL,
    [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [CreatedBy] [VARCHAR](100) NULL,
    [Modified] [DATETIME] NOT NULL DEFAULT GETDATE(),
    [ModifiedBy] [VARCHAR](100) NULL,
    CONSTRAINT [FK_CommunityActivities_Communities] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Communities]([Id]),
    CONSTRAINT [FK_CommunityActivities_ActivityTypes] FOREIGN KEY ([ActivityTypeId]) REFERENCES [dbo].[ActivityTypes]([Id])
)

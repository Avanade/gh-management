CREATE TABLE [dbo].[CommunityActivities]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [CommunityId] INT NOT NULL,
    [Date] DATETIME,
    [Name] VARCHAR(255) NOT NULL, 
    [ActivityTypeId] INT NOT NULL,
    [Url] VARCHAR(255) NULL,
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(100) NULL
    CONSTRAINT [FK_CommunityActivities_Communities] FOREIGN KEY (CommunityId) REFERENCES Communities(Id), 
    CONSTRAINT [FK_CommunityActivities_ActivityTypes] FOREIGN KEY (ActivityTypeId) REFERENCES ActivityTypes(Id)
)

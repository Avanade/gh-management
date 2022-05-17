CREATE TABLE [dbo].[CommunityActivities]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [CommunityId] INT NOT NULL, 
    [Username] VARCHAR(100) NOT NULL, 
    [Description] VARCHAR(255) NOT NULL, 
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(50) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(50) NULL
    CONSTRAINT [FK_CommunityActivities_Communities] FOREIGN KEY (CommunityId) REFERENCES Communities(Id), 
    CONSTRAINT [FK_CommunityActivities_Users] FOREIGN KEY (Username) REFERENCES Users(Username)
)

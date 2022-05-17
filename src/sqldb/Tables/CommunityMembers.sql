CREATE TABLE [dbo].[CommunityMembers]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [CommunityId] INT NOT NULL, 
    [Username] VARCHAR(100) NULL, 
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(50) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(50) NULL
    CONSTRAINT [FK_CommunityMembers_Communities] FOREIGN KEY (CommunityId) REFERENCES Communities(Id), 
    CONSTRAINT [FK_CommunityMembers_Users] FOREIGN KEY (Username) REFERENCES Users(Username)
)

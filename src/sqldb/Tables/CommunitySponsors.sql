CREATE TABLE [dbo].[CommunitySponsors]
(
  [Id] INT NOT NULL PRIMARY KEY IDENTITY,
  [CommunityId] INT NOT NULL,
  [UserPrincipalName] VARCHAR(100) NOT NULL,
  [Created] DATETIME NOT NULL DEFAULT getdate(), 
  [CreatedBy] VARCHAR(100) NULL, 
  [Modified] DATETIME NOT NULL DEFAULT getdate(), 
  [ModifiedBy] VARCHAR(100) NULL
  CONSTRAINT FK_CommunitySponsors_Communities FOREIGN KEY (CommunityId) REFERENCES Communities(Id),
  CONSTRAINT FK_CommunitySponsors_Users FOREIGN KEY (UserPrincipalName) REFERENCES Users(UserPrincipalName),
  INDEX IX_CommunityId_UserPrincipalName UNIQUE (CommunityId ASC, UserPrincipalName ASC)
)

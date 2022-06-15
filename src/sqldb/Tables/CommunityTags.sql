CREATE TABLE [dbo].[CommunityTags]
(
  [Id] INT NOT NULL PRIMARY KEY IDENTITY,
  [CommunityId] INT NOT NULL,
  [Tag] VARCHAR(20) NOT NULL
  CONSTRAINT FK_CommunityTags_Communities FOREIGN KEY (CommunityId) REFERENCES Communities(Id),
  INDEX IX_CommunityId_Tag UNIQUE (CommunityId ASC, Tag ASC)
)

CREATE TABLE [dbo].[CommunityTag] (
  [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
  [CommunityId] [INT] NOT NULL,
  [Tag] [VARCHAR](20) NOT NULL,
  CONSTRAINT [FK_CommunityTag_Community] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Community]([Id]),
  CONSTRAINT [AK_CommunityId_Tag] UNIQUE ([CommunityId], [Tag])
)

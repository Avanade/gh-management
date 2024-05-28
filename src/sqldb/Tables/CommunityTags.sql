CREATE TABLE [dbo].[CommunityTags] (
  [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
  [CommunityId] [INT] NOT NULL,
  [Tag] [VARCHAR](20) NOT NULL,
  CONSTRAINT [FK_CommunityTags_Communities] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Communities]([Id]),
  CONSTRAINT [AK_CommunityId_Tag] UNIQUE ([CommunityId], [Tag])
)

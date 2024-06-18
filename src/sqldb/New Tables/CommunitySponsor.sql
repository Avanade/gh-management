CREATE TABLE [dbo].[CommunitySponsor] (
  [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
  [CommunityId] [INT] NOT NULL,
  [UserPrincipalName] [VARCHAR](100) NOT NULL,
  [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
  [CreatedBy] [VARCHAR](100) NULL,
  [Modified] [DATETIME] NOT NULL DEFAULT GETDATE(),
  [ModifiedBy] [VARCHAR](100) NULL,
  CONSTRAINT [FK_CommunitySponsor_Community] FOREIGN KEY ([CommunityId]) REFERENCES [dbo].[Community]([Id]),
  CONSTRAINT [FK_CommunitySponsor_User] FOREIGN KEY ([UserPrincipalName]) REFERENCES [dbo].[User]([UserPrincipalName]),
  CONSTRAINT [AK_CommunityId_UserPrincipalName] UNIQUE ([CommunityId], [UserPrincipalName])
)

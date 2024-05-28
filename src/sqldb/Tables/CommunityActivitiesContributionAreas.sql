CREATE TABLE [dbo].[CommunityActivitiesContributionAreas] (
  [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
  [CommunityActivityId] [INT] NOT NULL,
  [ContributionAreaId] [INT] NOT NULL,
  [IsPrimary] [BIT] NOT NULL DEFAULT 0,
  [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
  [CreatedBy] [VARCHAR](100) NULL,
  [Modified] [DATETIME] NOT NULL DEFAULT GETDATE(),
  [ModifiedBy] [VARCHAR](100) NULL,
  CONSTRAINT [FK_CommunityActivitiesCA_CommunityActivity] FOREIGN KEY ([CommunityActivityId]) REFERENCES [dbo].[CommunityActivities]([Id]),
  CONSTRAINT [FK_CommunityActivitiesCA_ContributionAreas] FOREIGN KEY ([ContributionAreaId]) REFERENCES [dbo].[ContributionAreas]([Id])
)

CREATE TABLE [dbo].[CommunityActivityContributionArea] (
  [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
  [CommunityActivityId] [INT] NOT NULL,
  [ContributionAreaId] [INT] NOT NULL,
  [IsPrimary] [BIT] NOT NULL DEFAULT 0,
  CONSTRAINT [FK_CommunityActivityContributionArea_CommunityActivity] FOREIGN KEY ([CommunityActivityId]) REFERENCES [dbo].[CommunityActivity]([Id]),
  CONSTRAINT [FK_CommunityActivityContributionArea_ContributionArea] FOREIGN KEY ([ContributionAreaId]) REFERENCES [dbo].[ContributionArea]([Id])
)

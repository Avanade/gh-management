CREATE TABLE [dbo].[CommunityActivitiesContributionAreas]
(
  [Id] INT NOT NULL PRIMARY KEY IDENTITY,
  [CommunityActivityId] INT NOT NULL,
  [ContributionAreaId] INT NOT NULL,
  [IsPrimary] BIT NOT NULL DEFAULT 0,
  [Created] DATETIME NOT NULL DEFAULT getdate(), 
  [CreatedBy] VARCHAR(100) NULL, 
  [Modified] DATETIME NOT NULL DEFAULT getdate(), 
  [ModifiedBy] VARCHAR(100) NULL
  CONSTRAINT FK_CommunityActivitiesCA_CommunityActivity FOREIGN KEY (CommunityActivityId) REFERENCES CommunityActivities(Id),
  CONSTRAINT FK_CommunityActivitiesCA_ContributionAreas FOREIGN KEY (ContributionAreaId) REFERENCES ContributionAreas(Id)
)

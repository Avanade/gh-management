CREATE TABLE [dbo].[ContributionArea] (
  [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
  [Name] [VARCHAR](100) NOT NULL,
  [Created] [DATETIME] NOT NULL DEFAULT GETDATE(),
  [CreatedBy] [VARCHAR](100) NULL,
  [Modified] [DATETIME] NULL,
  [ModifiedBy] [VARCHAR](100) NULL
)

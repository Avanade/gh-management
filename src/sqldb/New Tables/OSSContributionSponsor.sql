CREATE TABLE [dbo].[OSSContributionSponsor] (
    [Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
    [Name] [VARCHAR](50) NOT NULL,
    [IsArchived] [BIT] NOT NULL DEFAULT 0
)

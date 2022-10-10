CREATE TABLE [dbo].[CommunityActivitiesHelpTypes](
	[Id] [int] NOT NULL PRIMARY KEY IDENTITY,
	[CommunityActivityId] [int] NOT NULL,
	[HelpTypeId] [int] NOT NULL,
	[Details] [nchar](100) NOT NULL
)

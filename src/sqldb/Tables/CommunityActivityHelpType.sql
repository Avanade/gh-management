CREATE TABLE [dbo].[CommunityActivityHelpType] (
	[Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
	[CommunityActivityId] [INT] NOT NULL,
	[HelpTypeId] [INT] NOT NULL,
	[Details] [NCHAR](100) NOT NULL
)

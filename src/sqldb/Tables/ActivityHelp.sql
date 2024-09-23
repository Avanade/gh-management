CREATE TABLE [dbo].[ActivityHelp] (
	[ActivityId] [INT] NOT NULL,
	[HelpTypeId] [INT] NOT NULL,
	[Details] [VARCHAR](255),
	CONSTRAINT [FK_ActivityHelp_HelpType] FOREIGN KEY ([HelpTypeId]) REFERENCES [dbo].[HelpType]([Id]),
	CONSTRAINT [FK_ActivityHelp_CommunityActivity] FOREIGN KEY ([ActivityId]) REFERENCES [dbo].[CommunityActivity]([Id])
)
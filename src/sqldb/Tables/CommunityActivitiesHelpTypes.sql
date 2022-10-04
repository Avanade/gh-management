/****** Object:  Table [dbo].[CommunityActivitiesHelpTypes]    Script Date: 15/07/2022 8:44:55 am ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[CommunityActivitiesHelpTypes](
	[Id] [int] NOT NULL PRIMARY KEY IDENTITY,,
	[CommunityActivityId] [int] NOT NULL,
	[HelpTypeId] [int] NOT NULL,
	[Details] [nchar](100) NOT NULL
) ON [PRIMARY]
GO


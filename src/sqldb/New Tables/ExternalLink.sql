CREATE TABLE [dbo].[ExternalLink] (
	[Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
	[IconSVG] [VARCHAR](1000) NOT NULL,
	[Hyperlink] [VARCHAR](100) NOT NULL,
	[LinkName] [VARCHAR](100) NOT NULL,
	[IsEnabled] [BIT] NOT NULL,
	[Created] [DATETIME] NOT NULL,
	[CreatedBy] [VARCHAR](100) NOT NULL,
	[Modified] [DATETIME] NULL,
	[ModifiedBy] [VARCHAR](100) NULL
)

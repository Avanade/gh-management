CREATE TABLE [dbo].[ExternalLinks] (
	[Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
	[IconSVG] [VARCHAR](1000) NULL,
	[Hyperlink] [VARCHAR](100) NULL,
	[LinkName] [VARCHAR](100) NULL,
	[Enabled] [BIT] NULL,
	[Created] [DATETIME] NOT NULL,
	[CreatedBy] [VARCHAR](100) NULL,
	[Modified] [DATETIME] NOT NULL,
	[ModifiedBy] [VARCHAR](100) NULL
)

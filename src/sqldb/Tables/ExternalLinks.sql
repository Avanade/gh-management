
CREATE TABLE [dbo].[ExternalLinks]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY,
	[IconSVG] [varchar](1000) NULL,
	[Hyperlink] [varchar](100) NULL,
	[LinkName] [varchar](100) NULL,
	[Enabled] [bit] NULL,
	[Created] [datetime] NOT NULL,
	[CreatedBy] [varchar](100) NULL,
	[Modified] [datetime] NOT NULL,
	[ModifiedBy] [varchar](100) NULL
)
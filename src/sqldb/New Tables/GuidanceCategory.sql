CREATE TABLE [dbo].[GuidanceCategory] (
	[Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
	[Name] [VARCHAR](100) NULL,
	[Created] [DATETIME] NULL,
	[CreatedBy] [VARCHAR](50) NULL,
	[Modified] [DATETIME] NULL,
	[ModifiedBy] [VARCHAR](50) NULL
)
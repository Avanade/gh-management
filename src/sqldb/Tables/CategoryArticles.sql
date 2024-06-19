CREATE TABLE [dbo].[CategoryArticles] (
	[Id] [INT] NOT NULL PRIMARY KEY IDENTITY,
	[Name] [VARCHAR](100) NULL,
	[Url] [VARCHAR](255) NULL,
	[Body] [VARCHAR](2000) NULL,
	[CategoryId] [INT] NULL,
	[Created] [DATETIME] NULL,
	[CreatedBy] [VARCHAR](50) NULL,
	[Modified] [DATETIME] NULL,
	[ModifiedBy] [VARCHAR](50) NULL,
	CONSTRAINT [FK_CategoryArticles_Category] FOREIGN KEY([CategoryId]) REFERENCES [dbo].[Category]([Id])
)

 
CREATE TABLE [dbo].[CategoryArticles](
	[Id] [int] IDENTITY(1,1) NOT NULL,
	[Name] [varchar](100) NULL,
	[Url] [varchar](100) NULL,
	[Body] [varchar](2000) NULL,
	[CategoryId] [int] NULL,
	[Created] [datetime] NULL,
	[CreatedBy] [varchar](50) NULL,
	[Modified] [datetime] NULL,
	[ModifiedBy] [varchar](50) NULL,
 CONSTRAINT [PK_CategoryArticles] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

ALTER TABLE [dbo].[CategoryArticles]  WITH CHECK ADD  CONSTRAINT [FK_CategoryArticles_Category] FOREIGN KEY([CategoryId])
REFERENCES [dbo].[Category] ([Id])
GO

ALTER TABLE [dbo].[CategoryArticles] CHECK CONSTRAINT [FK_CategoryArticles_Category]
GO



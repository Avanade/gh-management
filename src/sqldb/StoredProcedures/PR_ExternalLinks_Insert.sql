CREATE PROCEDURE PR_ExternalLinks_Insert
	@SVGName VARCHAR(100),
	@IconSVG VARCHAR(1000),
	@Hyperlink VARCHAR(100),
	@LinkName VARCHAR(100),
	@Category VARCHAR(100),
	@Enabled VARCHAR(100),
	@CreatedBy VARCHAR(100)

AS
BEGIN

INSERT INTO [dbo].[ExternalLinks] ( 
			[SVGName],
			[IconSVG],
			[Hyperlink],
			[LinkName],
			[Category],
			[Enabled],
            [Created],
            [CreatedBy],
            [Modified],
            [ModifiedBy]
)
VALUES
           ( 
			@SVGName,
			@IconSVG,
			@Hyperlink,
			@LinkName,
			@Category,
			@Enabled,
            GETDATE(),
            @CreatedBy,
            GETDATE(),
            @CreatedBy
            )
END
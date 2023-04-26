CREATE PROCEDURE PR_ExternalLinks_Insert
	@SVGName VARCHAR(100),
	@IconSVG VARCHAR(1000),
	@Category VARCHAR(100),
	@CreatedBy VARCHAR(100),
	@Enabled VARCHAR(100)
AS
BEGIN

INSERT INTO [dbo].[ExternalLinks] ( 
			[SVGName],
			[IconSVG],
			[Category],
            [Created],
            [CreatedBy],
            [Modified],
            [ModifiedBy],
			[Enabled]
)
VALUES
           ( 
			@SVGName,
			@IconSVG,
			@Category,
            GETDATE(),
            @CreatedBy,
            GETDATE(),
            @CreatedBy,
			@Enabled
            )
END
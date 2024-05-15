CREATE PROCEDURE PR_ExternalLinks_Insert
	@IconSVG VARCHAR(100),
	@Hyperlink VARCHAR(100),
	@LinkName VARCHAR(100),
	@Enabled VARCHAR(100),
	@CreatedBy VARCHAR(100)

AS
BEGIN

INSERT INTO [dbo].[ExternalLinks] ( 
			[IconSVG],
			[Hyperlink],
			[LinkName],
			[Enabled],
            [Created],
            [CreatedBy],
            [Modified],
            [ModifiedBy]
)
VALUES
           ( 
			@IconSVG,
			@Hyperlink,
			@LinkName,
			@Enabled,
            GETDATE(),
            @CreatedBy,
            GETDATE(),
            @CreatedBy
            )
END
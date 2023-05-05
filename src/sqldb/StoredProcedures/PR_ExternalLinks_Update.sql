CREATE PROCEDURE PR_ExternalLinks_Update
(
			@Id int,
			@SVGName VARCHAR(100),
			@Hyperlink VARCHAR(100),
			@LinkName VARCHAR(100),
			@Category VARCHAR(100),
			@Enabled VARCHAR(100),
			@CreatedBy VARCHAR(100),
			@ModifiedBy VARCHAR(100)
) AS
BEGIN

UPDATE [dbo].[ExternalLinks]
SET   [SVGName] = @SVGName, 
	  [Hyperlink] = @Hyperlink,
	  [LinkName] = @LinkName,
	  [Category] = @Category,
	  [Enabled] = @Enabled,
      [Created] = GETDATE(),
      [CreatedBy] =  @CreatedBy,
      [Modified] = GETDATE(),
	  [ModifiedBy] = @ModifiedBy
WHERE  [Id] = @Id

END
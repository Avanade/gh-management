CREATE PROCEDURE PR_ExternalLinks_Update
(
			@Id int,
			@IconSVG VARCHAR(100),
			@Hyperlink VARCHAR(100),
			@LinkName VARCHAR(100),
			@Enabled VARCHAR(100),
			@ModifiedBy VARCHAR(100)
) AS
BEGIN

UPDATE [dbo].[ExternalLinks]

SET   [IconSVG] = @IconSVG, 
	  [Hyperlink] = @Hyperlink,
	  [LinkName] = @LinkName,
	  [Enabled] = @Enabled,
      [Modified] = GETDATE(),
	  [ModifiedBy] = @ModifiedBy

WHERE  [Id] = @Id

END
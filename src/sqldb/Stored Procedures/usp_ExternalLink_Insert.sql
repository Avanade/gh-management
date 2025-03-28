CREATE PROCEDURE [dbo].[usp_ExternalLink_Insert]
	@IconSVG [VARCHAR](100),
	@Hyperlink [VARCHAR](100),
	@LinkName [VARCHAR](100),
	@IsEnabled [VARCHAR](100),
	@CreatedBy [VARCHAR](100)
AS
BEGIN
  INSERT INTO [dbo].[ExternalLink] ( 
    [IconSVG],
    [Hyperlink],
    [LinkName],
    [IsEnabled],
    [Created],
    [CreatedBy]
  )
  OUTPUT
    [INSERTED].[Id],
    [INSERTED].[Created]
  VALUES (
    @IconSVG,
    @Hyperlink,
    @LinkName,
    @IsEnabled,
    GETDATE(),
    @CreatedBy
  );
END

CREATE PROCEDURE [dbo].[usp_ExternalLink_Update]
  @Id [INT],
  @IconSVG [VARCHAR](100),
  @Hyperlink [VARCHAR](100),
  @LinkName [VARCHAR](100),
  @IsEnabled [VARCHAR](100),
  @ModifiedBy [VARCHAR](100)
AS
BEGIN
  UPDATE [dbo].[ExternalLink]
  SET
    [IconSVG] = @IconSVG, 
    [Hyperlink] = @Hyperlink,
    [LinkName] = @LinkName,
    [IsEnabled] = @IsEnabled,
    [Modified] = GETDATE(),
    [ModifiedBy] = @ModifiedBy
  OUTPUT
    [INSERTED].[Modified]
  WHERE  [Id] = @Id
END

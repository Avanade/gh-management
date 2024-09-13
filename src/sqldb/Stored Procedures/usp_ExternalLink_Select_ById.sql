CREATE PROCEDURE [dbo].[usp_ExternalLink_Select_ById]
  @Id [INT]
AS
BEGIN
  SELECT 
    [Id],
    [IconSVG],
    [Hyperlink],
    [LinkName],
    [IsEnabled],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy]
  FROM [dbo].[ExternalLink]
  WHERE [Id] = @Id
END

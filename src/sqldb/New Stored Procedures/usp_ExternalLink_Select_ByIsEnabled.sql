CREATE PROCEDURE [dbo].[usp_ExternalLink_Select_ByIsEnabled]
	@Enabled [BIT] = true
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
  WHERE [IsEnabled] = @Enabled
  ORDER BY [Id] DESC
END

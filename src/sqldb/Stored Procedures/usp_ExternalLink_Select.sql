CREATE PROCEDURE [dbo].[usp_ExternalLink_Select]
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
  FROM	[dbo].[ExternalLink]
  ORDER BY [ID] DESC
END

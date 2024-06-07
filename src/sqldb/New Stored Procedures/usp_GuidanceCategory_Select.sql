CREATE PROCEDURE [dbo].[usp_GuidanceCategory_Select]
AS 
BEGIN
  SELECT
    [Id],
    [Name],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy]
  FROM [dbo].[GuidanceCategory]
END

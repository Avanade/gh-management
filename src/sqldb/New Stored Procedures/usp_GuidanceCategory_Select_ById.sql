CREATE PROCEDURE [dbo].[usp_GuidanceCategory_Select_ById]
  @Id [INT]
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
  WHERE [Id] = @Id
END

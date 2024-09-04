CREATE PROCEDURE [dbo].[usp_ContributionArea_Select_ById]
  @Id [INT]
AS
BEGIN
  SET NOCOUNT ON

  SELECT
    [Id],
    [Name],
    [Created],
    [CreatedBy],
    [Modified],
    [ModifiedBy]
  FROM [dbo].[ContributionArea]
  WHERE [Id] = @Id
END
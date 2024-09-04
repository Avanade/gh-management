CREATE PROCEDURE [dbo].[usp_ContributionArea_Select]
AS
BEGIN
    SELECT
      [Id],
      [Name],
      [Created],
      [CreatedBy],
      [Modified],
      [ModifiedBy]
    FROM [dbo].[ContributionArea]
END
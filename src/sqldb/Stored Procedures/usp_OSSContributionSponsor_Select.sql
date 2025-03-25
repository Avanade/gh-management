CREATE PROCEDURE [dbo].[usp_OSSContributionSponsor_Select]
AS
BEGIN
  SELECT 
    [Id],
    [Name],
    [IsArchived]
  FROM [dbo].[OSSContributionSponsor]
  ORDER BY [Id] ASC
END

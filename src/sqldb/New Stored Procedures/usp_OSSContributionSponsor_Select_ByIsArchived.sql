CREATE PROCEDURE [dbo].[usp_OSSContributionSponsor_Select_ByIsArchived]
    @IsArchived [BIT] = 0
AS
BEGIN
  SELECT 
    [Id],
    [Name],
    [IsArchived]
  FROM [dbo].[OSSContributionSponsor]
  WHERE [IsArchived]=@IsArchived
  ORDER BY [Id] ASC

END
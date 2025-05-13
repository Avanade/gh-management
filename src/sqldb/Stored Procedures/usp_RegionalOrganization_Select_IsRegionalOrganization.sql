CREATE PROCEDURE [dbo].[usp_RegionalOrganization_Select_IsRegionalOrganization]
  @IsEnabled [BIT] = NULL -- NULL = ALL | 1 = ENABLED | 0 = DISABLED
AS
BEGIN
    SELECT
      [Id],
      [Name],
      [IsEnabled],
      [Created],
      [CreatedBy],
      [Modified],
      [ModifiedBy]
    FROM 
      [dbo].[RegionalOrganization]
    WHERE
      (
        [IsRegionalOrganization] = 1
      ) AND
      (@IsEnabled IS NULL 
      OR
      (
        @IsEnabled IS NOT NULL AND
        IsEnabled = @IsEnabled
      ))
END
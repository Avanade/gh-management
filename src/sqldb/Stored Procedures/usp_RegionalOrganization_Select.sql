CREATE PROCEDURE [dbo].[usp_RegionalOrganization_Select]
  @IsEnabled [BIT] = NULL -- NULL = ALL | 1 = ENABLED | 0 = DISABLED
AS
BEGIN
    SELECT
      [Id],
      [Name],
      [IsRegionalOrganization],
      [IsIndexRepoEnabled],
      [IsCopilotRequestEnabled],
      [IsAccessRequestEnabled],
      [IsEnabled],
      [Created],
      [CreatedBy],
      [Modified],
      [ModifiedBy]
    FROM 
      [dbo].[RegionalOrganization]
    WHERE
      @IsEnabled IS NULL 
      OR
      (
        @IsEnabled IS NOT NULL AND
        IsEnabled = @IsEnabled
      )
END
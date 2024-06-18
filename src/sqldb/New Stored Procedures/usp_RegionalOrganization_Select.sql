CREATE PROCEDURE [dbo].[usp_RegionalOrganization_Select]
AS
BEGIN
    SELECT
      [Id],
      [Name]
    FROM [dbo].[RegionalOrganization] 

END
CREATE PROCEDURE [dbo].[usp_RegionalOrganization_Select_ById]
  @Id [INT]
AS
BEGIN
    SELECT 
      [Id],
      [Name]
    FROM [dbo].[RegionalOrganization] 
    WHERE Id = @Id
END

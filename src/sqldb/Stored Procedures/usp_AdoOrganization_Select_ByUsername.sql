CREATE PROCEDURE [dbo].[usp_AdoOrganization_Select_ByUsername]
  @Username [VARCHAR](100)
AS
BEGIN
  SELECT 
    [AO].[Id],
    [AO].[Name],
    [AO].[Purpose],
    [AO].[Created],
    [AO].[CreatedBy]
  FROM [dbo].[AdoOrganization] AS [AO]
  WHERE [AO].[CreatedBy] = @Username
  ORDER BY [AO].[Created] DESC
END

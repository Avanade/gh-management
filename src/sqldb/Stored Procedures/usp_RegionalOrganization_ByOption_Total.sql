CREATE PROCEDURE [dbo].[usp_RegionalOrganization_Select_ByOption_Total]
	@Search [VARCHAR](100) = ''
AS
BEGIN
    SELECT
      COUNT(*) AS [Total]
    FROM 
      [dbo].[RegionalOrganization] AS [RO]
    INNER JOIN 
      STRING_SPLIT(@Search, ' ') AS [SS] ON ([RO].[Name] LIKE '%'+[SS].[value]+'%')
END
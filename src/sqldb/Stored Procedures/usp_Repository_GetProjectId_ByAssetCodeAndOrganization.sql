CREATE PROCEDURE [dbo].[usp_Repository_GetProjectId_ByAssetCodeAndOrganization]
	@AssetCode [VARCHAR](50),
	@Organization [VARCHAR](100)
AS
BEGIN
  SELECT 
    [Id]
  FROM 
    [dbo].[Repository] 
  WHERE 
    [Organization] = @Organization AND [AssetCode] = @AssetCode
END
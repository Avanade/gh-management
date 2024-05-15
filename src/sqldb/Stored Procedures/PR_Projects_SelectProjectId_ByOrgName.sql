CREATE PROCEDURE [dbo].[PR_Projects_SelectProjectId_ByOrgName]
(
	@AssetCode VARCHAR(50),
	@Organization VARCHAR(100)
)
AS
BEGIN
    SELECT 
        Id
    FROM 
        [dbo].[Projects] 
    WHERE 
        Organization = @Organization AND AssetCode = @AssetCode
END

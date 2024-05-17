CREATE PROCEDURE [dbo].[PR_Projects_Select_ByOrgName]
(
	@AssetCode VARCHAR(50),
	@Organization VARCHAR(100)
)
AS
BEGIN
    SELECT 
        *
    FROM 
        [dbo].[Projects] 
    WHERE 
        Organization = @Organization AND AssetCode = @AssetCode
END
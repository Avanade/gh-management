CREATE PROCEDURE [dbo].[PR_Projects_SelectProjectId_ByOrgName]
(
	@Name VARCHAR(50),
	@Organization VARCHAR(100)
)
AS
BEGIN
    SELECT 
        Id
    FROM 
        [dbo].[Projects] 
    WHERE 
        Organization = @Organization AND Name = @Name
END
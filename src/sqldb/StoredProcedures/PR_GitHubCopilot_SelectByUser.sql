CREATE PROCEDURE [dbo].[PR_GitHubCopilot_SelectByUser]
    @Username VARCHAR(100)
AS
BEGIN
    SELECT 
        GC.Id,
        RO.Name,
        GC.GitHubUsername,
        GC.Created
    FROM [dbo].[GitHubCopilot] GC
    LEFT JOIN RegionalOrganizations RO ON GC.Region = RO.Id
    WHERE GC.CreatedBy=@Username
    ORDER BY GC.Created DESC
END
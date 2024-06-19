CREATE PROCEDURE [dbo].[PR_Organizations_SelectByUser]
    @Username VARCHAR(100)
AS
BEGIN
    SELECT 
        O.Id,
        RO.Name,
        O.ClientName,
        O.ProjectName,
        O.WBS,
        A.Name,
        O.Created
    FROM [dbo].[Organizations] O
    LEFT JOIN RegionalOrganizations RO ON O.Region = RO.Id
    LEFT JOIN ApprovalStatus A ON A.Id = O.ApprovalStatusId
    WHERE O.CreatedBy=@Username
    ORDER BY O.Created DESC
END
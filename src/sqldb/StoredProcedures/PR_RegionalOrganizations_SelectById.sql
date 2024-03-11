CREATE PROCEDURE [dbo].[PR_RegionalOrganizations_SelectById]
(
 @Id int
 )
AS
BEGIN

    SELECT * 
    FROM [dbo].[RegionalOrganizations] 
    WHERE Id = @Id

END
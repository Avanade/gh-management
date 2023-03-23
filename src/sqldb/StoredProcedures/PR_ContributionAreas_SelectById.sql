CREATE PROCEDURE [dbo].[PR_ContributionAreas_SelectById]
(
	@Id INT
)
AS
BEGIN
    SET NOCOUNT ON
    SELECT 
	* 
    FROM ContributionAreas
	WHERE
		Id = @Id 
END
CREATE PROCEDURE PR_Applications_Select_ById

	@Id UNIQUEIDENTIFIER

AS

SELECT
Id, [Name]
FROM Applications
WHERE
IsActive = 1
AND Id = @Id
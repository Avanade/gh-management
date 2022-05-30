CREATE PROCEDURE PR_Applications_Select_ById

	@Id UNIQUEIDENTIFIER

AS

SELECT
dbo.UidToString(Id) [Id], [Name]
FROM Applications
WHERE
IsActive = 1
AND Id = @Id
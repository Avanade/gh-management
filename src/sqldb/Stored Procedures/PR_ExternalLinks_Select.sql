CREATE PROCEDURE PR_ExternalLinks_Select

AS
BEGIN

SELECT * FROM	[dbo].[ExternalLinks]

ORDER BY id desc

END
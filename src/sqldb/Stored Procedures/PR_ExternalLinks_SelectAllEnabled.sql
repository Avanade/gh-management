CREATE PROCEDURE PR_ExternalLinks_SelectAllEnabled
	@Enabled as bit = true

AS
BEGIN

SELECT * FROM	[dbo].[ExternalLinks]
WHERE 
	Enabled = @Enabled

ORDER BY Id desc

END
CREATE PROCEDURE PR_ExternalLinks_Select
	@CreatedBy as varchar(100)
AS
BEGIN

SELECT * FROM	[dbo].[ExternalLinks]

WHERE 
	CreatedBy = @CreatedBy
ORDER BY id desc

END

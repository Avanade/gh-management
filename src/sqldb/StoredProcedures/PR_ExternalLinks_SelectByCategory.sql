CREATE PROCEDURE PR_ExternalLinks_SelectByCategory
	@CreatedBy as varchar(100),
	@Category as VARCHAR(100)

AS
BEGIN

SELECT * FROM	[dbo].[ExternalLinks]

WHERE 
	CreatedBy = @CreatedBy
	and
	Category = @Category

END
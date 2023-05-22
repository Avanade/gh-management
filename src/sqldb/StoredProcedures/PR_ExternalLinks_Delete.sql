CREATE PROCEDURE PR_ExternalLinks_Delete
@Id int
AS
DELETE FROM ExternalLinks

where Id = @Id
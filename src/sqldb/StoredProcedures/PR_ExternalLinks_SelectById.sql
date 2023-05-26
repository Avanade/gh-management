CREATE PROCEDURE PR_ExternalLinks_SelectById
  @id int 
AS
BEGIN
  SELECT * FROM
	ExternalLinks
  WHERE  
    Id = @id
END
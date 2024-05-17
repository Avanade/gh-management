CREATE PROCEDURE PR_GitHubAccess_SelectByObjectID
  @ObjectId VARCHAR(100) 
AS
BEGIN
  SELECT * FROM [dbo].[GitHubAccess]
  WHERE  
    [ObjectId] = @ObjectId
END
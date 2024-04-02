CREATE PROCEDURE PR_Communities_InitCommunityType
 
AS
BEGIN
    UPDATE [dbo].[Communities]
    SET   
      CommunityType =  IIF(IsExternal=1, 'external', 'internal')
    WHERE CommunityType IS NULL
END
GO
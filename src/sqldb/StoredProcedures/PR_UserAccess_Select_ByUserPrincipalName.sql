CREATE PROCEDURE PR_UserAccess_Select_ByUserPrincipalName
(
	@UserPrincipalName varchar(100)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
 
SELECT [Id],
       [ProjectId],
       [UserPrincipalName],
       [Created],
       [CreatedBy],
       [Modified],
       [ModifiedBy]
  FROM 
       [dbo].[UserAccess]
  WHERE  
       [UserPrincipalName] = @UserPrincipalName
  ORDER BY [Created] DESC
END
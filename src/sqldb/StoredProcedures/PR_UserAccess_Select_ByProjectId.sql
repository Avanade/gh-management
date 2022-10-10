CREATE PROCEDURE PR_UserAccess_Select_ByProjectId
(
	@ProjectId INT
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
       [ProjectId] = @ProjectId
END
CREATE PROCEDURE PR_UserAccess_Select

AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
 

    -- Insert statements for procedure here
SELECT [Id],
       [ProjectId],
       [UserPrincipalName],
       [Created],
       [CreatedBy],
       [Modified],
       [ModifiedBy]
  FROM [dbo].[UserAccess]
END
CREATE PROCEDURE PR_UserAccess_Select
(
	@Id Int
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
 

    -- Insert statements for procedure here
SELECT [Id]
      ,[ProjectId]
      ,[Username]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[UserAccess]
  WHERE  [Id] =@Id or @Id IS NULL
END
CREATE PROCEDURE [dbo].[PR_Users_Select]
(
	@Username varchar(100)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here

SELECT [Username]
      ,[FirstName]
      ,[LastName]
      ,[Email]
      ,[Created]
      ,[CreatedBy]
      ,[Modified]
      ,[ModifiedBy]
  FROM [dbo].[Users]
  where  
	  [Username] =@Username or @Username IS NULL

END
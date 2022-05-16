Create PROCEDURE [dbo].[PR_Users_Update]
(
	@Username varchar(100),
    @FirstName varchar(50),
    @LastName varchar(50),
    @Email varchar(50),
    @Created datetime,
    @CreatedBy varchar(50),
    @Modified datetime,
    @ModifiedBy varchar(50)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here
UPDATE [dbo].[Users]
   SET [Username] = @Username
      ,[FirstName] = @FirstName
      ,[LastName] = @LastName
      ,[Email] = @Email
      ,[Created] = @Created
      ,[CreatedBy] = @CreatedBy
      ,[Modified] = @Modified
      ,[ModifiedBy] = @ModifiedBy
 WHERE  
	   [Username] =@Username or @Username IS NULL
END
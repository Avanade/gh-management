CREATE PROCEDURE PR_Users_Insert
(
			@Username varchar(100)
           ,@FirstName varchar(50)
           ,@LastName varchar(50)
           ,@Email varchar(50)
           ,@Created datetime
           ,@CreatedBy varchar(50)
           ,@Modified datetime
           ,@ModifiedBy varchar(50)
)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

    -- Insert statements for procedure here

INSERT INTO [dbo].[Users]
           ([Username]
           ,[FirstName]
           ,[LastName]
           ,[Email]
           ,[Created]
           ,[CreatedBy]
           ,[Modified]
           ,[ModifiedBy])
     VALUES
           (@Username
           ,@FirstName
           ,@LastName
           ,@Email
           ,GETDATE()
           ,@CreatedBy
           ,GETDATE()
           ,@ModifiedBy)
END
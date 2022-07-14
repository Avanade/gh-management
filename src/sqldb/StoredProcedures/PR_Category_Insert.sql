
create PROCEDURE  [dbo].[PR_Category_Insert]
(
	@Name varchar(50),
	@CreatedBy  varchar(50),
	@ModifiedBy  varchar(50) ,
	@Id  int =null
)
AS
BEGIN   
    SET NOCOUNT ON 
	DECLARE @returnID AS INT
	   select @id= id from  [Category] where name = @Name
	  
	IF (@Id=0  OR @Id is null )
	BEGIN
			INSERT INTO [dbo].[Category]
					   ([Name] ,
					   [Created],
					   [CreatedBy],
					   [Modified],
					   [ModifiedBy]
					   )
				 VALUES
					   (@Name  
					   ,GETDATE()
					   ,@CreatedBy
					   ,GETDATE()
					   ,@ModifiedBy)

			SET @returnID = SCOPE_IDENTITY()
			SELECT @returnID Id
	end
	ELSE 
	begin 
	 EXEC	[dbo].[PR_Category_Update] @Id , @Name  ,@CreatedBy ,@ModifiedBy

	SELECT @Id Id
	end
END

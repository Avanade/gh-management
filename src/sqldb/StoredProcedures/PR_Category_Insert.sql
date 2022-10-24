
CREATE PROCEDURE  [dbo].[PR_Category_Insert]
(
	@Name VARCHAR(50),
	@CreatedBy  VARCHAR(50),
	@ModifiedBy  VARCHAR(50) ,
	@Id  INT = NULL
)
AS
BEGIN   
    SET NOCOUNT ON 
	DECLARE @returnID AS INT
	   SELECT @id= id FROM  [Category] WHERE NAME = @Name
	  
	IF (@Id=0  OR @Id IS NULL )
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
	END
	ELSE 
	BEGIN 
	 EXEC	[dbo].[PR_Category_Update] @Id , @Name  ,@CreatedBy ,@ModifiedBy

	SELECT @Id Id
	END
END

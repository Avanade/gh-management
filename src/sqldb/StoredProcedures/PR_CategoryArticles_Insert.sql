
CREATE PROCEDURE [dbo].[PR_CategoryArticles_Insert]
(
			@Name VARCHAR(100),
			@Url VARCHAR(100),
			@Body VARCHAR(2000),
			@CategoryId INT,
            @CreatedBy VARCHAR(50),
            @ModifiedBy VARCHAR(50),
			@Id  INT = NULL
) AS
BEGIN
	DECLARE @returnID AS INT
 
	--IF NOT EXISTS (SELECT Id FROM [Communities] WHERE id  = @Id  )
	IF (@Id=0  )
	BEGIN
 

		INSERT INTO [dbo].[CategoryArticles]
           ([Name]
           ,[Url]
           ,[Body]
           ,[CategoryId]
           ,[Created]
           ,[CreatedBy]
           ,[Modified]
           ,[ModifiedBy])
     VALUES
           (@Name
           ,@Url
           ,@Body
           ,@CategoryId
           ,GETDATE()
           ,@CreatedBy
           ,GETDATE()
           ,@ModifiedBy)
			 SET @returnID = SCOPE_IDENTITY()


 				SELECT @returnID Id
	END
	ELSE
	BEGIN
	EXEC	  [dbo].[PR_CategoryArticles_Update] 	
			@Id   ,
			@Name  ,
			@Url ,
			@Body  ,
			@CategoryId  ,
            @CreatedBy ,
            @ModifiedBy  
			

	SELECT @Id Id
	END
END
 



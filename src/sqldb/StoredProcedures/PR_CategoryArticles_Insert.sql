
CREATE PROCEDURE [dbo].[PR_CategoryArticles_Insert]
(
			@Name varchar(100),
			@Url varchar(100),
			@Body varchar(2000),
			@CategoryId int,
            @CreatedBy varchar(50),
            @ModifiedBy varchar(50),
			@Id  int =null
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
	end
	else
	begin
	EXEC	  [dbo].[PR_CategoryArticles_Update] 	
			@Id   ,
			@Name  ,
			@Url ,
			@Body  ,
			@CategoryId  ,
            @CreatedBy ,
            @ModifiedBy  
			

	SELECT @Id Id
	end
end
 



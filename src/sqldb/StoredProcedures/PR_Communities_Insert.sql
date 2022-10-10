
CREATE PROCEDURE [dbo].[PR_Communities_Insert]
(
			@Name VARCHAR(50),
			@Url VARCHAR(255),
			@Description VARCHAR(255),
			@Notes VARCHAR(255),
			@TradeAssocId VARCHAR(255),
			@IsExternal INT,
			@CreatedBy  VARCHAR(50),
			@ModifiedBy  VARCHAR(50) ,
			@Id  INT = NULL
) AS
BEGIN
	DECLARE @returnID AS INT
 
	--IF NOT EXISTS (SELECT Id FROM [Communities] WHERE id  = @Id  )
	IF (@Id=0  )
	BEGIN
 

			INSERT INTO [dbo].[Communities]
					   ([Name]
					   ,[Url]
					   ,[Description]
					   ,[Notes]
					   ,[TradeAssocId]
					   ,[IsExternal]
					   ,[Created]
					   ,[CreatedBy]
					   ,[Modified]
					   ,[ModifiedBy])
				 VALUES
					   (@Name
					   ,@Url
					   ,@Description
					   ,@Notes
					   ,@TradeAssocId
					   ,@IsExternal
					   ,GETDATE()
					   ,@CreatedBy
					   ,GETDATE()
					   ,@ModifiedBy	)
			 SET @returnID = SCOPE_IDENTITY()


 				SELECT @returnID Id
	END
	ELSE
	BEGIN
	EXEC	  [dbo].[PR_Communities_Update]
		@Id ,
		@Name ,
		@Url ,
		@Description ,
		@Notes ,
		@TradeAssocId ,
		@IsExternal,
		@CreatedBy ,
		@ModifiedBy

	SELECT @Id Id
	END
END

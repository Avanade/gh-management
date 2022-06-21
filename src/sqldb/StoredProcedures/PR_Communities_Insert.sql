
create PROCEDURE [dbo].[PR_Communities_Insert]
(
			@Name varchar(50),
			@Url varchar(255),
			@Description varchar(255),
			@Notes varchar(255),
			@TradeAssocId varchar(255),
			@IsExternal int,
			@CreatedBy  varchar(50),
			@ModifiedBy  varchar(50) ,
			@Id  int =null
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
	end
	else
	begin
	EXEC	  [dbo].[PR_Communities_Update]
		@Id ,
		@Name ,
		@Url ,
		@Description ,
		@Notes ,
		@TradeAssocId ,
		@CreatedBy ,
		@ModifiedBy

	SELECT @Id Id
	end
end

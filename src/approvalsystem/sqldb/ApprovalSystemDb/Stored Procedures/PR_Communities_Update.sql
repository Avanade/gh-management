/****** Object:  StoredProcedure [dbo].[PR_Projects_Insert]    Script Date: 5/31/2022 9:32:57 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
create PROCEDURE [dbo].[PR_Communities_Update]
(
			@Name varchar(50),
			@Url varchar(255),
			@Description varchar(255),
			@Notes varchar(255),
			@TradeAssocId varchar(255),
			@CreatedBy  varchar(50),
			@ModifiedBy  varchar(50)
) AS
BEGIN
UPDATE [dbo].[Communities]
   SET [Name] = @Name
      ,[Url] = @Url
      ,[Description] = @Description
      ,[Notes] = @Notes
      ,[TradeAssocId] = @TradeAssocId
      ,[Created] =GETDATE()
      ,[CreatedBy] = @CreatedBy
      ,[Modified] = GETDATE()
      ,[ModifiedBy] = @ModifiedBy
 WHERE  [Name] = @Name


end

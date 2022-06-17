/****** Object:  StoredProcedure [dbo].[PR_Projects_Insert]    Script Date: 5/31/2022 9:32:57 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
create PROCEDURE [dbo].[PR_Communities_Insert]
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
	DECLARE @Id AS INT
            
            INSERT INTO [dbo].[Communities]
                  ([Name]
                  ,[Url]
                  ,[Description]
                  ,[Notes]
                  ,[TradeAssocId]
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
                  ,GETDATE()
                  ,@CreatedBy
                  ,GETDATE()
                  ,@ModifiedBy	)
            SET @Id = SCOPE_IDENTITY()
 	SELECT @Id Id

end

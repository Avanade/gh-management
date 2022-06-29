create PROCEDURE [dbo].[PR_Communities_Update]
(
			@Id int,
			@Name varchar(50),
			@Url varchar(255),
			@Description varchar(255),
			@Notes varchar(255),
			@TradeAssocId varchar(255),
            @IsExternal int,
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
      ,IsExternal=@IsExternal
      ,[Created] =GETDATE()
      ,[CreatedBy] = @CreatedBy
      ,[Modified] = GETDATE()
      ,[ModifiedBy] = @ModifiedBy
 WHERE  [Id] = @Id

 delete from CommunitySponsors where CommunityId = @Id
  delete from CommunityTags where CommunityId = @Id
end

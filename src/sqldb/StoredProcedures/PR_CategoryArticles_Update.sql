create PROCEDURE [dbo].[PR_CategoryArticles_Update]
(
			@Id int,
            @Name varchar(100),
			@Url varchar(100),
			@Body varchar(2000),
			@CategoryId int,
            @CreatedBy varchar(50),
            @ModifiedBy varchar(50)
) AS
BEGIN
 
UPDATE [dbo].[CategoryArticles]
   SET [Name] = @Name
      ,[Url] =  @Url
      ,[Body] = @Body
      ,[CategoryId] = @CategoryId
      ,[Created] = GETDATE()
      ,[CreatedBy] = @CreatedBy
      ,[Modified] =  GETDATE()
      ,[ModifiedBy] =@ModifiedBy
 WHERE  [Id] = @Id
 
 
end

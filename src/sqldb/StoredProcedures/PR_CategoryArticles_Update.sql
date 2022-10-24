CREATE PROCEDURE [dbo].[PR_CategoryArticles_Update]
(
			@Id INT,
            @Name VARCHAR(100),
			@Url VARCHAR(100),
			@Body VARCHAR(2000),
			@CategoryId INT,
            @CreatedBy VARCHAR(50),
            @ModifiedBy VARCHAR(50)
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
END
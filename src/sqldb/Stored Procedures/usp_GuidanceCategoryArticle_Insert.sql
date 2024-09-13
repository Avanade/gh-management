CREATE PROCEDURE [dbo].[usp_GuidanceCategoryArticle_Insert]
  @Name [VARCHAR](100),
  @Url [VARCHAR](255),
  @Body [VARCHAR](2000),
  @CategoryId [INT],
  @CreatedBy [VARCHAR](50),
  @ModifiedBy [VARCHAR](50),
  @Id  [INT] = NULL
AS
BEGIN
	IF NOT EXISTS ( SELECT * FROM [dbo].[GuidanceCategoryArticle] WHERE [Id] = @Id )
	  BEGIN
  		INSERT INTO [dbo].[GuidanceCategoryArticle]
      (
        [Name],
        [Url],
        [Body],
        [GuidanceCategoryId],
        [Created],
        [CreatedBy],
        [Modified],
        [ModifiedBy]
      )
      VALUES
      (
        @Name,
        @Url,
        @Body,
        @CategoryId,
        GETDATE(),
        @CreatedBy,
        GETDATE(),
        @ModifiedBy
      )
			SET @Id = SCOPE_IDENTITY()
	  END
	ELSE
	  BEGIN
	    EXEC [dbo].[usp_GuidanceCategoryArticle_Update] 	
			  @Id,
			  @Name,
			  @Url,
			  @Body,
			  @CategoryId,
        @CreatedBy,
        @ModifiedBy  
  	END
  SELECT @Id as [Id]
END

CREATE PROCEDURE [dbo].[usp_GuidanceCategoryArticle_Insert]
(
  @Name [VARCHAR](100),
  @Url [VARCHAR](255),
  @Body [VARCHAR](2000),
  @CategoryId [INT],
  @CreatedBy [VARCHAR](50),
  @ModifiedBy [VARCHAR](50),
  @Id  [INT] = NULL
)
AS
BEGIN
	DECLARE @returnID AS [INT]

	IF ( @Id = 0 )
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
			SET @returnID = SCOPE_IDENTITY()

      SELECT @returnID AS [Id]
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
	    SELECT @Id AS [Id]
  	END
END

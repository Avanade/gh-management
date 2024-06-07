CREATE PROCEDURE [dbo].[usp_GuidanceCategoryArticle_Update]
(
  @Id INT,
  @Name [VARCHAR](100),
  @Url [VARCHAR](255),
  @Body [VARCHAR](2000),
  @CategoryId [INT],
  @CreatedBy [VARCHAR](50),
  @ModifiedBy [VARCHAR](50)
)
AS
BEGIN
  UPDATE [dbo].[GuidanceCategoryArticle]
  SET
    [Name] = @Name,
    [Url] =  @Url,
    [Body] = @Body,
    [GuidanceCategoryId] = @CategoryId,
    [Created] = GETDATE(),
    [CreatedBy] = @CreatedBy,
    [Modified] =  GETDATE(),
    [ModifiedBy] =@ModifiedBy
  WHERE [Id] = @Id
END
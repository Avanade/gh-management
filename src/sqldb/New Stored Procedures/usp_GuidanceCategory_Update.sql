CREATE PROCEDURE [dbo].[usp_GuidanceCategory_Update]
(
  @Id [INT],
  @Name [VARCHAR](50),	
  @CreatedBy [VARCHAR](50),
  @ModifiedBy [VARCHAR](50) 
)
AS
BEGIN
  UPDATE [dbo].[GuidanceCategory]
  SET
    [Name] = @Name,
    [Created] = GETDATE(),
    [CreatedBy] = @CreatedBy,
    [Modified] = GETDATE(),
    [ModifiedBy] = @ModifiedBy
  WHERE [Id] = @Id 
END

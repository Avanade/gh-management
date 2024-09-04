CREATE PROCEDURE [dbo].[usp_CommunityActivity_Insert]
  @CommunityId [INT],
  @Name [VARCHAR](255),
  @ActivityTypeId [INT],
  @Url [VARCHAR](255),
  @Date [DATE],
  @CreatedBy [VARCHAR](50)
AS
BEGIN
  DECLARE @Id AS [INT]
  INSERT INTO [dbo].[CommunityActivity]
    (
    [CommunityId],
    [Name],
    [ActivityTypeId],
    [Url],
    [Date],
    [Created],
    [CreatedBy]
    )
  VALUES
    (
      @CommunityId,
      @Name,
      @ActivityTypeId,
      @Url,
      @Date,
      GETDATE(),
      @CreatedBy
    )
  SET @Id = SCOPE_IDENTITY()
  SELECT @Id AS [Id]
END
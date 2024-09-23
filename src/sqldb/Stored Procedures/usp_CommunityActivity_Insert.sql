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
  OUTPUT
    [INSERTED].[Id],
    [INSERTED].[Created]
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
END
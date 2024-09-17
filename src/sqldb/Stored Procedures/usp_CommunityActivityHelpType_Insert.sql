CREATE PROCEDURE [dbo].[usp_CommunityActivityHelpType_Insert]
  @ActivityId [INT],
  @HelpTypeId [INT],
  @Details [VARCHAR](100)
AS
BEGIN
  DECLARE @Id AS [INT]

  INSERT INTO [dbo].[CommunityActivityHelpType]
  (
    [CommunityActivityId],
    [HelpTypeId],
    [Details]
  )
  OUTPUT
    [INSERTED].[Id]
  VALUES
  (
    @ActivityId,
    @HelpTypeId,
    @Details
  )

  SET @Id = SCOPE_IDENTITY()
  SELECT @Id [Id]
END
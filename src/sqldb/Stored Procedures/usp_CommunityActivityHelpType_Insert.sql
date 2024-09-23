CREATE PROCEDURE [dbo].[usp_CommunityActivityHelpType_Insert]
  @ActivityId [INT],
  @HelpTypeId [INT],
  @Details [VARCHAR](100)
AS
BEGIN
  INSERT INTO [dbo].[ActivityHelp]
  (
    [ActivityId],
    [HelpTypeId],
    [Details]
  )
  VALUES
  (
    @ActivityId,
    @HelpTypeId,
    @Details
  )
END
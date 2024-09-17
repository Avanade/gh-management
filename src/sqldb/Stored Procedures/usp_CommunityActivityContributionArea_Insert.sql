CREATE PROCEDURE [dbo].[usp_CommunityActivityContributionArea_Insert]
	@CommunityActivityId [INT],
	@ContributionAreaId [INT],
	@IsPrimary [BIT]
AS
BEGIN
	INSERT INTO [dbo].[CommunityActivityContributionArea](
		[CommunityActivityId],
		[ContributionAreaId],
		[IsPrimary]
  	)
	OUTPUT
		[INSERTED].[Id]
	VALUES (
		@CommunityActivityId,
		@ContributionAreaId,
		@IsPrimary
  	)
END
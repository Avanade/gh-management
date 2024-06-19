CREATE PROCEDURE [dbo].[usp_CommunityActivityContributionArea_Insert]
	@CommunityActivityId [INT],
	@ContributionAreaId [INT],
	@IsPrimary [BIT],
	@CreatedBy [VARCHAR](50)
AS
BEGIN
	DECLARE @Id AS [INT]
  INSERT INTO [dbo].[CommunityActivityContributionArea](
		[CommunityActivityId],
		[ContributionAreaId],
		[IsPrimary],
		[Created],
		[CreatedBy]
  ) VALUES (
		@CommunityActivityId,
		@ContributionAreaId,
		@IsPrimary,
		GETDATE(),
		@CreatedBy
  )
	SET @Id = SCOPE_IDENTITY()
	SELECT @Id AS [Id]
END
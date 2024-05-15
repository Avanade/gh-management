CREATE PROCEDURE [dbo].[PR_CommunityActivitiesContributionAreas_Insert]
(
	@CommunityActivityId INT,
	@ContributionAreaId INT,
	@IsPrimary BIT,
	@CreatedBy VARCHAR(50)
)
AS
BEGIN
	DECLARE @Id AS INT
    INSERT INTO CommunityActivitiesContributionAreas(
		CommunityActivityId,
		ContributionAreaId,
		IsPrimary,
		Created,
		CreatedBy
    ) VALUES (
		@CommunityActivityId,
		@ContributionAreaId,
		@IsPrimary,
		GETDATE(),
		@CreatedBy
    )
	SET @Id = SCOPE_IDENTITY()
	SELECT @Id Id
END
GO


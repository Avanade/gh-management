/****** Object:  StoredProcedure [dbo].[PR_CommunityActivitiesContributionAreas_Insert]    Script Date: 6/16/2022 3:42:50 PM ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE PROCEDURE [dbo].[PR_CommunityActivitiesContributionAreas_Insert]
(
	@CommunityActivityId int,
	@ContributionAreaId int,
	@IsPrimary bit,
	@CreatedBy varchar(50)
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


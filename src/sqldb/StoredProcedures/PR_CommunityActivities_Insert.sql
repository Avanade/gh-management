/****** Object:  StoredProcedure [dbo].[PR_CommunitiesActivities_Insert]    Script Date: 6/16/2022 3:41:49 PM ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE PROCEDURE [dbo].[PR_CommunitiesActivities_Insert]
(
    @CommunityId int,
    @Name varchar(255),
    @ActivityTypeId int,
    @Url varchar(255),
	@Date date,
    @CreatedBy varchar(50)
)
AS
BEGIN
	DECLARE @Id AS INT
    INSERT INTO CommunityActivities(
        [CommunityId],
        [Name],
        [ActivityTypeId],
        [Url],
		[Date],
        [Created],
        [CreatedBy]
    ) VALUES (
        @CommunityId,
        @Name,
        @ActivityTypeId,
        @Url,
		@Date,
		GETDATE(),
        @CreatedBy
    )
	SET @Id = SCOPE_IDENTITY()
	SELECT @Id Id
END
GO
CREATE PROCEDURE [dbo].[PR_CommunityActivities_Insert]
(
    @CommunityId INT,
    @Name VARCHAR(255),
    @ActivityTypeId INT,
    @Url VARCHAR(255),
	@Date DATE,
    @CreatedBy VARCHAR(50)
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
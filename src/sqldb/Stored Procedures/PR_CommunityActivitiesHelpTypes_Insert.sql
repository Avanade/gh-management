CREATE PROCEDURE [dbo].[PR_CommunityActivitiesHelpTypes_Insert] 
(
	@ActivityActivityId INT,
	@HelpTypeId INT,
	@Details VARCHAR(100)
)
AS
BEGIN
	DECLARE @Id AS INT

	INSERT INTO [dbo].[CommunityActivitiesHelpTypes]
			   ([CommunityActivityId]
			   ,[HelpTypeId]
			   ,[Details])
		 VALUES
			   (@ActivityActivityId
			   ,@HelpTypeId
			   ,@Details)
	
	SET @Id = SCOPE_IDENTITY()

	SELECT @Id Id
END
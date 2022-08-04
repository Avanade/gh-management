/****** Object:  StoredProcedure [dbo].[PR_CommunityActivitiesHelpTypes_Insert]    Script Date: 15/07/2022 9:12:12 am ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

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